package config

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/crossplane/upjet/v2/pkg/config"
)

// ---------------------------------------------------------------------------
// External name pattern helpers
//
// Snowflake's Terraform provider uses two ID encoding functions:
//
//   EncodeResourceIdentifier  — newer, used by stable resources.
//     AccountObjectIdentifier → Name() (bare, no quotes)
//     DatabaseObjectIdentifier → FullyQualifiedName() → "db"."name"
//     SchemaObjectIdentifier → FullyQualifiedName() → "db"."schema"."name"
//     Multiple parts join with |.
//
//   EncodeSnowflakeID  — legacy, still used by preview resources.
//     AccountObjectIdentifier → Name() (bare, no quotes)
//     DatabaseObjectIdentifier → db|name (unquoted, pipe-separated parts)
//     SchemaObjectIdentifier → db|schema|name (unquoted, pipe-separated parts)
//
// The helpers below encode the corresponding template for each pattern.
// See .work/snowflakedb/snowflake/pkg/helpers/ for the canonical implementation.
// ---------------------------------------------------------------------------

// SchemaObjectIdentifier returns external name config for Snowflake's
// 3-part SchemaObjectIdentifier fully-qualified-name format:
//
//	"db"."schema"."name"
//
// The name is omitted from spec.forProvider (set via metadata.name →
// crossplane.io/external-name annotation).
//
// Used by stable-family resources that use EncodeResourceIdentifier
// with a SchemaObjectIdentifier argument — the ID is produced by
// FullyQualifiedName() which quotes each segment.
func SchemaObjectIdentifier() config.ExternalName {
	return config.TemplatedStringAsIdentifier("name",
		`"{{ .parameters.database }}"."{{ .parameters.schema }}"."{{ .external_name }}"`)
}

// PipeSeparatedIdentifier returns external name config for Snowflake's
// pipe-separated 3-part ID format:
//
//	db|schema|name
//
// Used by preview-family resources that use EncodeSnowflakeID
// with a SchemaObjectIdentifier argument — the ID is produced by
// joining databaseName, schemaName, and Name() with |.
func PipeSeparatedIdentifier() config.ExternalName {
	return config.TemplatedStringAsIdentifier("name",
		`{{ .parameters.database }}|{{ .parameters.schema }}|{{ .external_name }}`)
}

// DatabaseSchemaIdentifier returns external name config for Snowflake's
// 2-part quoted ID format:
//
//	"db"."name"
//
// Used by resources whose ID is a DatabaseObjectIdentifier
// (e.g. schema, database_role).
func DatabaseSchemaIdentifier() config.ExternalName {
	return config.TemplatedStringAsIdentifier("name",
		`"{{ .parameters.database }}"."{{ .external_name }}"`)
}

// OneOfIdentifier returns a TemplatedStringAsIdentifier for resources whose
// ID is composed of the external name followed by a conditional chain of
// one-of-many mutually-exclusive parameter-prefix pairs.
//
// prefix is a Go template string that precedes the conditional chain
// (e.g. "{{ .external_name }}|" or '"{{ .external_name }}"|').
// params maps each TF parameter name to its ID prefix label.
//
// The generated template iterates over sorted keys for deterministic output.
// All params use else-if (no bare else fallback), safe because these
// parameters are ExactlyOneOf in the TF schema.
//
// Example:
//
//	OneOfIdentifier("database_role_name",
//	    `{{ .external_name }}|`,
//	    map[string]string{
//	        "parent_database_role_name": "DATABASE_ROLE",
//	        "parent_role_name":          "ROLE",
//	        "share_name":                "SHARE",
//	    },
//	)
//
// produces (sorted alpha by key):
//
//	{{ .external_name }}|{{ if .parameters.parent_database_role_name }}DATABASE_ROLE|"{{ .parameters.parent_database_role_name }}"{{ else if .parameters.parent_role_name }}ROLE|"{{ .parameters.parent_role_name }}"{{ else if .parameters.share_name }}SHARE|"{{ .parameters.share_name }}"{{ end }}
func OneOfIdentifier(nameField, prefix string, params map[string]string) config.ExternalName {
	keys := slices.Sorted(maps.Keys(params))

	var tmpl strings.Builder
	tmpl.WriteString(prefix)

	for i, k := range keys {
		if i == 0 {
			fmt.Fprintf(&tmpl, `{{ if .parameters.%s }}%s|"{{ .parameters.%s }}"`, k, params[k], k)
		} else {
			fmt.Fprintf(&tmpl, `{{ else if .parameters.%s }}%s|"{{ .parameters.%s }}"`, k, params[k], k)
		}
	}
	tmpl.WriteString(`{{ end }}`)

	return config.TemplatedStringAsIdentifier(nameField, tmpl.String())
}

// =============================================================================
// Shared helpers for grant resources with dynamic-length compound IDs.
// =============================================================================

// grantAllFutureParts handles the common OnAll/OnFuture sub-block pattern
// shared by grant_ownership, grant_privileges_to_account_role, and
// grant_privileges_to_database_role.
//
// The block has object_type_plural and either in_database or in_schema.
// Returns the suffix parts: [object_type_plural, InDatabase|InSchema, identifier].
func grantAllFutureParts(block map[string]any) ([]string, error) {
	objTypePlural, _ := block["object_type_plural"].(string)
	if objTypePlural == "" {
		return nil, fmt.Errorf("object_type_plural is required")
	}
	parts := []string{objTypePlural}
	if db, _ := block["in_database"].(string); db != "" {
		parts = append(parts, "InDatabase", db)
	} else if s, _ := block["in_schema"].(string); s != "" {
		parts = append(parts, "InSchema", s)
	} else {
		return nil, fmt.Errorf("requires in_database or in_schema")
	}
	return parts, nil
}

// onSchemaBlockSuffix constructs the ID suffix for an on_schema block.
// Returns (grant sub-type, suffix parts, found).
func onSchemaBlockSuffix(onSchema map[string]any) (string, []string, bool) {
	if v, _ := onSchema["schema_name"].(string); v != "" {
		return "OnSchema", []string{v}, true
	}
	if v, _ := onSchema["all_schemas_in_database"].(string); v != "" {
		return "OnAllSchemasInDatabase", []string{v}, true
	}
	if v, _ := onSchema["future_schemas_in_database"].(string); v != "" {
		return "OnFutureSchemasInDatabase", []string{v}, true
	}
	return "", nil, false
}

// onSchemaObjectBlockSuffix constructs the ID suffix for an on_schema_object block.
// Returns (grant sub-type, suffix parts, found, error).
// "all_privileges" is checked to determine whether to include an extra
// OnObject sub-type marker for the single-object case. This matches the
// TF provider's ID generation logic (from acceptance test ID comments).
func onSchemaObjectBlockSuffix(so map[string]any, allPrivileges bool) (string, []string, bool, error) {
	// Single object: object_type + object_name
	if objType, _ := so["object_type"].(string); objType != "" {
		objName, _ := so["object_name"].(string)
		if allPrivileges {
			return "OnObject", []string{objType, objName}, true, nil
		}
		// Without all_privileges, the TF provider omits the OnObject
		// marker in the ID: OnSchemaObject|<type>|<name>
		return "", []string{objType, objName}, true, nil
	}
	// OnAll: all[0] sub-block
	if allRaw, _ := so["all"].([]any); len(allRaw) > 0 {
		allBlock, ok := allRaw[0].(map[string]any)
		if !ok {
			return "", nil, false, fmt.Errorf("invalid 'all' block")
		}
		parts, err := grantAllFutureParts(allBlock)
		return "OnAll", parts, true, err
	}
	// OnFuture: future[0] sub-block
	if futureRaw, _ := so["future"].([]any); len(futureRaw) > 0 {
		futureBlock, ok := futureRaw[0].(map[string]any)
		if !ok {
			return "", nil, false, fmt.Errorf("invalid 'future' block")
		}
		parts, err := grantAllFutureParts(futureBlock)
		return "OnFuture", parts, true, err
	}
	return "", nil, false, nil
}

// grantPrivilegesBaseStr constructs the common base prefix for grant_privileges
// resources: <role_name>|<with_grant_option>|<always_apply>|<privileges>.
// privileges is formatted as comma-separated sorted list, or "ALL" when
// all_privileges is true.
func grantPrivilegesBaseStr(roleName string, parameters map[string]any) string {
	wgo := "false"
	if v, ok := parameters["with_grant_option"]; ok {
		switch b := v.(type) {
		case bool:
			if b {
				wgo = "true"
			}
		case string:
			if b == "true" {
				wgo = "true"
			}
		}
	}
	aa := "false"
	if v, ok := parameters["always_apply"]; ok {
		switch b := v.(type) {
		case bool:
			if b {
				aa = "true"
			}
		case string:
			if b == "true" {
				aa = "true"
			}
		}
	}

	var privs string
	if allPriv, ok := parameters["all_privileges"]; ok {
		switch b := allPriv.(type) {
		case bool:
			if b {
				privs = "ALL"
			}
		case string:
			if b == "true" {
				privs = "ALL"
			}
		}
	}
	if privs == "" {
		if privRaw, ok := parameters["privileges"].([]any); ok && len(privRaw) > 0 {
			p := make([]string, len(privRaw))
			for i, v := range privRaw {
				p[i] = fmt.Sprint(v)
			}
			slices.Sort(p)
			privs = strings.Join(p, ",")
		}
	}

	return fmt.Sprintf("%s|%s|%s|%s", roleName, wgo, aa, privs)
}

// GrantOwnershipIdentifier returns an ExternalName for snowflake_grant_ownership.
//
// The TF resource ID is a compound, variable-length pipe-separated string:
//
//	<role_type>|<role_identifier>|<outbound_privileges>|<grant_type>|<grant_data>
//
// where:
// - role_type is "ToAccountRole" or "ToDatabaseRole"
// - role_identifier is the fully qualified name of the role
// - outbound_privileges is "COPY", "REVOKE", or empty
// - grant_type is "OnObject", "OnAll", or "OnFuture"
// - grant_data is structured based on grant_type:
//
// OnObject: <object_type>|<object_name>
// OnAll/OnFuture (InDatabase): <object_type_plural>|InDatabase|<database>
// OnAll/OnFuture (InSchema):   <object_type_plural>|InSchema|<schema>
//
// The ID is generated by the TF provider on Create and reconstructed from
// parameters for import.
func GrantOwnershipIdentifier() config.ExternalName {
	return config.NewExternalNameFrom(
		config.IdentifierFromProvider,
		config.WithGetIDFn(func(fn config.GetIDFn, ctx context.Context, externalName string, parameters map[string]any, providerConfig map[string]any) (string, error) {
			if id, err := buildGrantOwnershipID(parameters); err == nil && id != "" {
				return id, nil
			}
			return fn(ctx, externalName, parameters, providerConfig)
		}),
	)
}

// buildGrantOwnershipID constructs the import ID for snowflake_grant_ownership
// from the resource parameters.
func buildGrantOwnershipID(parameters map[string]any) (string, error) {
	// Determine role type and role ID from mutually exclusive params.
	var roleType, roleID string
	if v, _ := parameters["account_role_name"].(string); v != "" {
		roleType = "ToAccountRole"
		roleID = v
	} else if v, _ := parameters["database_role_name"].(string); v != "" {
		roleType = "ToDatabaseRole"
		roleID = v
	}
	if roleType == "" {
		return "", fmt.Errorf("grant_ownership: neither account_role_name nor database_role_name is set")
	}

	outboundPrivileges, _ := parameters["outbound_privileges"].(string)
	parts := []string{roleType, roleID, outboundPrivileges}

	// "on" is Block List, Min: 1, Max: 1.
	onRaw, _ := parameters["on"].([]any)
	if len(onRaw) == 0 {
		return "", fmt.Errorf("grant_ownership: 'on' block is required")
	}
	onBlock, ok := onRaw[0].(map[string]any)
	if !ok {
		return "", fmt.Errorf("grant_ownership: invalid 'on' block")
	}

	// OnObject: object_type + object_name
	if objType, _ := onBlock["object_type"].(string); objType != "" {
		objName, _ := onBlock["object_name"].(string)
		parts = append(parts, "OnObject", objType, objName)
		return strings.Join(parts, "|"), nil
	}

	// OnAll: all[0] sub-block (object_type_plural + in_database or in_schema)
	if allRaw, _ := onBlock["all"].([]any); len(allRaw) > 0 {
		allBlock, ok := allRaw[0].(map[string]any)
		if !ok {
			return "", fmt.Errorf("grant_ownership: invalid 'all' block")
		}
		suffix, err := grantAllFutureParts(allBlock)
		if err != nil {
			return "", fmt.Errorf("grant_ownership: 'all' block: %w", err)
		}
		parts = append(parts, "OnAll")
		parts = append(parts, suffix...)
		return strings.Join(parts, "|"), nil
	}

	// OnFuture: future[0] sub-block (object_type_plural + in_database or in_schema)
	if futureRaw, _ := onBlock["future"].([]any); len(futureRaw) > 0 {
		futureBlock, ok := futureRaw[0].(map[string]any)
		if !ok {
			return "", fmt.Errorf("grant_ownership: invalid 'future' block")
		}
		suffix, err := grantAllFutureParts(futureBlock)
		if err != nil {
			return "", fmt.Errorf("grant_ownership: 'future' block: %w", err)
		}
		parts = append(parts, "OnFuture")
		parts = append(parts, suffix...)
		return strings.Join(parts, "|"), nil
	}

	return "", fmt.Errorf("grant_ownership: unable to determine grant type from 'on' block")
}

// =============================================================================
// Snowflake-specific external name identifier functions
// =============================================================================

// GrantPrivilegesToAccountRoleIdentifier returns an ExternalName for
// snowflake_grant_privileges_to_account_role.
//
// The TF resource ID is a compound, variable-length pipe-separated string:
//
//	<account_role_name>|<with_grant_option>|<always_apply>|<privileges>|<grant_type>|<grant_data>
func GrantPrivilegesToAccountRoleIdentifier() config.ExternalName {
	return config.NewExternalNameFrom(
		config.IdentifierFromProvider,
		config.WithGetIDFn(func(fn config.GetIDFn, ctx context.Context, externalName string, parameters map[string]any, providerConfig map[string]any) (string, error) {
			if id, err := buildGrantPrivilegesToAccountRoleID(parameters); err == nil && id != "" {
				return id, nil
			}
			return fn(ctx, externalName, parameters, providerConfig)
		}),
	)
}

func buildGrantPrivilegesToAccountRoleID(parameters map[string]any) (string, error) {
	roleName, _ := parameters["account_role_name"].(string)
	if roleName == "" {
		return "", fmt.Errorf("grant_privileges_to_account_role: account_role_name is required")
	}

	allPrivileges := false
	if v, ok := parameters["all_privileges"]; ok {
		switch b := v.(type) {
		case bool:
			if b {
				allPrivileges = true
			}
		case string:
			if b == "true" {
				allPrivileges = true
			}
		}
	}

	base := grantPrivilegesBaseStr(roleName, parameters)

	// OnAccount: boolean field
	if v, ok := parameters["on_account"]; ok {
		switch b := v.(type) {
		case bool:
			if b {
				return base + "|OnAccount", nil
			}
		case string:
			if b == "true" {
				return base + "|OnAccount", nil
			}
		}
	}

	// OnAccountObject: block with object_type + object_name
	if oaRaw, _ := parameters["on_account_object"].([]any); len(oaRaw) > 0 {
		oaBlock, ok := oaRaw[0].(map[string]any)
		if !ok {
			return "", fmt.Errorf("grant_privileges_to_account_role: invalid 'on_account_object' block")
		}
		objType, _ := oaBlock["object_type"].(string)
		objName, _ := oaBlock["object_name"].(string)
		return fmt.Sprintf("%s|OnAccountObject|%s|%s", base, objType, objName), nil
	}

	// OnSchema: block with one of schema_name, all_schemas_in_database, future_schemas_in_database
	if schemaRaw, _ := parameters["on_schema"].([]any); len(schemaRaw) > 0 {
		schemaBlock, ok := schemaRaw[0].(map[string]any)
		if !ok {
			return "", fmt.Errorf("grant_privileges_to_account_role: invalid 'on_schema' block")
		}
		subType, suffix, found := onSchemaBlockSuffix(schemaBlock)
		if !found {
			return "", fmt.Errorf("grant_privileges_to_account_role: on_schema block variant not recognized")
		}
		return fmt.Sprintf("%s|OnSchema|%s|%s", base, subType, strings.Join(suffix, "|")), nil
	}

	// OnSchemaObject: block with single object, OnAll, or OnFuture
	if soRaw, _ := parameters["on_schema_object"].([]any); len(soRaw) > 0 {
		soBlock, ok := soRaw[0].(map[string]any)
		if !ok {
			return "", fmt.Errorf("grant_privileges_to_account_role: invalid 'on_schema_object' block")
		}
		subType, suffix, found, err := onSchemaObjectBlockSuffix(soBlock, allPrivileges)
		if err != nil {
			return "", fmt.Errorf("grant_privileges_to_account_role: %w", err)
		}
		if !found {
			return "", fmt.Errorf("grant_privileges_to_account_role: on_schema_object block variant not recognized")
		}
		if subType == "" {
			// all_privileges=false, single object: omit OnObject sub-type marker
			return fmt.Sprintf("%s|OnSchemaObject|%s", base, strings.Join(suffix, "|")), nil
		}
		return fmt.Sprintf("%s|OnSchemaObject|%s|%s", base, subType, strings.Join(suffix, "|")), nil
	}

	return "", fmt.Errorf("grant_privileges_to_account_role: one of on_account, on_account_object, on_schema, or on_schema_object is required")
}

// GrantPrivilegesToDatabaseRoleIdentifier returns an ExternalName for
// snowflake_grant_privileges_to_database_role.
//
// The TF resource ID is a compound, variable-length pipe-separated string:
//
//	<database_role_name>|<with_grant_option>|<always_apply>|<privileges>|<grant_type>|<grant_data>
func GrantPrivilegesToDatabaseRoleIdentifier() config.ExternalName {
	return config.NewExternalNameFrom(
		config.IdentifierFromProvider,
		config.WithGetIDFn(func(fn config.GetIDFn, ctx context.Context, externalName string, parameters map[string]any, providerConfig map[string]any) (string, error) {
			if id, err := buildGrantPrivilegesToDatabaseRoleID(parameters); err == nil && id != "" {
				return id, nil
			}
			return fn(ctx, externalName, parameters, providerConfig)
		}),
	)
}

func buildGrantPrivilegesToDatabaseRoleID(parameters map[string]any) (string, error) {
	roleName, _ := parameters["database_role_name"].(string)
	if roleName == "" {
		return "", fmt.Errorf("grant_privileges_to_database_role: database_role_name is required")
	}

	allPrivileges := false
	if v, ok := parameters["all_privileges"]; ok {
		switch b := v.(type) {
		case bool:
			if b {
				allPrivileges = true
			}
		case string:
			if b == "true" {
				allPrivileges = true
			}
		}
	}

	base := grantPrivilegesBaseStr(roleName, parameters)

	// OnDatabase: string field
	if db, _ := parameters["on_database"].(string); db != "" {
		return fmt.Sprintf("%s|OnDatabase|%s", base, db), nil
	}

	// OnSchema: block with one of schema_name, all_schemas_in_database, future_schemas_in_database
	if schemaRaw, _ := parameters["on_schema"].([]any); len(schemaRaw) > 0 {
		schemaBlock, ok := schemaRaw[0].(map[string]any)
		if !ok {
			return "", fmt.Errorf("grant_privileges_to_database_role: invalid 'on_schema' block")
		}
		subType, suffix, found := onSchemaBlockSuffix(schemaBlock)
		if !found {
			return "", fmt.Errorf("grant_privileges_to_database_role: on_schema block variant not recognized")
		}
		return fmt.Sprintf("%s|OnSchema|%s|%s", base, subType, strings.Join(suffix, "|")), nil
	}

	// OnSchemaObject: block with single object, OnAll, or OnFuture
	if soRaw, _ := parameters["on_schema_object"].([]any); len(soRaw) > 0 {
		soBlock, ok := soRaw[0].(map[string]any)
		if !ok {
			return "", fmt.Errorf("grant_privileges_to_database_role: invalid 'on_schema_object' block")
		}
		subType, suffix, found, err := onSchemaObjectBlockSuffix(soBlock, allPrivileges)
		if err != nil {
			return "", fmt.Errorf("grant_privileges_to_database_role: %w", err)
		}
		if !found {
			return "", fmt.Errorf("grant_privileges_to_database_role: on_schema_object block variant not recognized")
		}
		if subType == "" {
			return fmt.Sprintf("%s|OnSchemaObject|%s", base, strings.Join(suffix, "|")), nil
		}
		return fmt.Sprintf("%s|OnSchemaObject|%s|%s", base, subType, strings.Join(suffix, "|")), nil
	}

	return "", fmt.Errorf("grant_privileges_to_database_role: one of on_database, on_schema, or on_schema_object is required")
}

// ---------------------------------------------------------------------------
// External name configurations
// ---------------------------------------------------------------------------

var ExternalNameConfigs = map[string]config.ExternalName{
	// =========================================================================
	// Stable resources (ShortGroup: stable)
	// =========================================================================

	// snowflake_account: import with '"<org>"."<account>"' — compound, but
	// the external name is just the bare account name. Org comes from the
	// provider config. AccountObjectIdentifier via EncodeResourceIdentifier.
	"snowflake_account": config.NameAsIdentifier,

	// snowflake_account_parameter: import with '<key>' — ID is the parameter
	// key. No 'name' field; uses 'key' via ParameterAsIdentifier.
	"snowflake_account_parameter": config.ParameterAsIdentifier("key"),

	// snowflake_account_role: import with '"<account_role_name>"' — bare name.
	"snowflake_account_role": config.NameAsIdentifier,

	// snowflake_account_session_policy_attachment: import format is the
	// fully qualified session policy name (FQN). The FQN is stored in the
	// session_policy_name field which doubles as the identifier.
	"snowflake_account_session_policy_attachment": config.ParameterAsIdentifier("session_policy_name"),

	// snowflake_api_authentication_integration_*: integration names.
	"snowflake_api_authentication_integration_with_authorization_code_grant": config.NameAsIdentifier,
	"snowflake_api_authentication_integration_with_client_credentials":       config.NameAsIdentifier,
	"snowflake_api_authentication_integration_with_jwt_bearer":               config.NameAsIdentifier,

	// snowflake_authentication_policy: SchemaObjectIdentifier.
	"snowflake_authentication_policy": SchemaObjectIdentifier(),

	// snowflake_catalog_integration_*: integration names.
	"snowflake_catalog_integration_aws_glue":       config.NameAsIdentifier,
	"snowflake_catalog_integration_iceberg_rest":   config.NameAsIdentifier,
	"snowflake_catalog_integration_object_storage": config.NameAsIdentifier,
	"snowflake_catalog_integration_open_catalog":   config.NameAsIdentifier,

	// snowflake_compute_pool: bare name.
	"snowflake_compute_pool": config.NameAsIdentifier,

	// snowflake_current_account: singleton — the literal string "current_account".
	"snowflake_current_account": config.IdentifierFromProvider,

	// snowflake_current_organization_account: import with '"<name>"' — bare name.
	"snowflake_current_organization_account": config.NameAsIdentifier,

	// snowflake_database: bare name (AccountObjectIdentifier).
	"snowflake_database": config.NameAsIdentifier,

	// snowflake_database_role: DatabaseObjectIdentifier — "db"."role".
	"snowflake_database_role": DatabaseSchemaIdentifier(),

	// snowflake_execute: random UUID generated by the TF provider.
	"snowflake_execute": config.IdentifierFromProvider,

	// snowflake_external_oauth_integration: bare name.
	"snowflake_external_oauth_integration": config.NameAsIdentifier,

	// snowflake_external_volume: bare name.
	"snowflake_external_volume": config.NameAsIdentifier,

	// snowflake_git_repository: SchemaObjectIdentifier.
	"snowflake_git_repository": SchemaObjectIdentifier(),

	// snowflake_grant_account_role: compound with conditional grantee type.
	// ID: '<role_name>|ROLE|<parent_role_name>' or
	//     '<role_name>|USER|<user_name>'.
	"snowflake_grant_account_role": OneOfIdentifier(
		"role_name",
		`"{{ .external_name }}"|`,
		map[string]string{
			"parent_role_name": "ROLE",
			"user_name":        "USER",
		},
	),

	// snowflake_grant_application_role: compound with conditional grantee type.
	// ID: '<app_role_fqn>|ACCOUNT_ROLE|<parent_account_role>' or
	//     '<app_role_fqn>|APPLICATION|<application>'.
	"snowflake_grant_application_role": OneOfIdentifier(
		"application_role_name",
		`{{ .external_name }}|`,
		map[string]string{
			"parent_account_role_name": "ACCOUNT_ROLE",
			"application_name":         "APPLICATION",
		},
	),

	// snowflake_grant_database_role: compound with conditional grantee type.
	// ID: '<db_role_fqn>|ROLE|<parent_role>' or
	//     '<db_role_fqn>|DATABASE_ROLE|<parent_db_role>' or
	//     '<db_role_fqn>|SHARE|<share>'.
	"snowflake_grant_database_role": OneOfIdentifier(
		"database_role_name",
		`{{ .external_name }}|`,
		map[string]string{
			"parent_role_name":          "ROLE",
			"parent_database_role_name": "DATABASE_ROLE",
			"share_name":                "SHARE",
		},
	),

	// snowflake_grant_ownership: compound ID with 5-7+ variable parts
	// (<target_role>|<fqn>|<outbound>|<kind>|<object_data>).
	// Custom GetIDFn reconstructs the variable-length import ID from parameters.
	"snowflake_grant_ownership": GrantOwnershipIdentifier(),

	// snowflake_grant_privileges_to_account_role: compound ID with 5-9+
	// variable parts (<role>|<with_grant>|<always>|<privileges>|<kind>|<data>).
	// Custom GetIDFn reconstructs the variable-length import ID from parameters.
	"snowflake_grant_privileges_to_account_role": GrantPrivilegesToAccountRoleIdentifier(),

	// snowflake_grant_privileges_to_database_role: compound ID with 6-9+
	// variable parts (<db_role>|<with_grant>|<always>|<privileges>|<kind>|<data>).
	// Custom GetIDFn reconstructs the variable-length import ID from parameters.
	"snowflake_grant_privileges_to_database_role": GrantPrivilegesToDatabaseRoleIdentifier(),

	// snowflake_grant_privileges_to_share: compound ID with conditional target.
	// ID: '<share>|<privileges>|OnDatabase|<db>' or OnSchema/OnTable/etc.
	"snowflake_grant_privileges_to_share": OneOfIdentifier(
		"to_share",
		`"{{ .external_name }}"|{{ .parameters.privileges }}|`,
		map[string]string{
			"on_database":             "OnDatabase",
			"on_schema":               "OnSchema",
			"on_table":                "OnTable",
			"on_all_tables_in_schema": "OnAllTablesInSchema",
			"on_tag":                  "OnTag",
			"on_view":                 "OnView",
		},
	),

	// snowflake_image_repository: SchemaObjectIdentifier.
	"snowflake_image_repository": SchemaObjectIdentifier(),

	// snowflake_legacy_service_user: same as user — bare name.
	"snowflake_legacy_service_user": config.NameAsIdentifier,

	// snowflake_listing: bare name.
	"snowflake_listing": config.NameAsIdentifier,

	// snowflake_masking_policy: SchemaObjectIdentifier.
	"snowflake_masking_policy": SchemaObjectIdentifier(),

	// snowflake_network_policy: bare name.
	"snowflake_network_policy": config.NameAsIdentifier,

	// snowflake_network_rule: SchemaObjectIdentifier.
	"snowflake_network_rule": SchemaObjectIdentifier(),

	// snowflake_oauth_integration_*: integration names.
	"snowflake_oauth_integration_for_custom_clients":       config.NameAsIdentifier,
	"snowflake_oauth_integration_for_partner_applications": config.NameAsIdentifier,

	// snowflake_password_policy: SchemaObjectIdentifier.
	"snowflake_password_policy": SchemaObjectIdentifier(),

	// snowflake_primary_connection: bare name.
	"snowflake_primary_connection": config.NameAsIdentifier,

	// snowflake_resource_monitor: bare name.
	"snowflake_resource_monitor": config.NameAsIdentifier,

	// snowflake_row_access_policy: SchemaObjectIdentifier.
	"snowflake_row_access_policy": SchemaObjectIdentifier(),

	// snowflake_saml2_integration: bare name.
	"snowflake_saml2_integration": config.NameAsIdentifier,

	// snowflake_schema: DatabaseObjectIdentifier — "db"."schema".
	"snowflake_schema": DatabaseSchemaIdentifier(),

	// snowflake_scim_integration: bare name.
	"snowflake_scim_integration": config.NameAsIdentifier,

	// snowflake_secondary_connection: bare name.
	"snowflake_secondary_connection": config.NameAsIdentifier,

	// snowflake_secondary_database: bare name.
	"snowflake_secondary_database": config.NameAsIdentifier,

	// snowflake_secret_with_*: SchemaObjectIdentifier (all variants).
	"snowflake_secret_with_authorization_code_grant": SchemaObjectIdentifier(),
	"snowflake_secret_with_basic_authentication":     SchemaObjectIdentifier(),
	"snowflake_secret_with_client_credentials":       SchemaObjectIdentifier(),
	"snowflake_secret_with_generic_string":           SchemaObjectIdentifier(),

	// snowflake_service: SchemaObjectIdentifier.
	"snowflake_service": SchemaObjectIdentifier(),

	// snowflake_service_user: same as user — bare name.
	"snowflake_service_user": config.NameAsIdentifier,

	// snowflake_session_policy: SchemaObjectIdentifier.
	"snowflake_session_policy": SchemaObjectIdentifier(),

	// snowflake_shared_database: bare name.
	"snowflake_shared_database": config.NameAsIdentifier,

	// snowflake_stage_external_*: SchemaObjectIdentifier (all variants).
	"snowflake_stage_external_azure":         SchemaObjectIdentifier(),
	"snowflake_stage_external_gcs":           SchemaObjectIdentifier(),
	"snowflake_stage_external_s3":            SchemaObjectIdentifier(),
	"snowflake_stage_external_s3_compatible": SchemaObjectIdentifier(),

	// snowflake_stage_internal: SchemaObjectIdentifier.
	"snowflake_stage_internal": SchemaObjectIdentifier(),

	// snowflake_storage_integration_*: integration names (provider-specific).
	"snowflake_storage_integration_aws":   config.NameAsIdentifier,
	"snowflake_storage_integration_azure": config.NameAsIdentifier,
	"snowflake_storage_integration_gcs":   config.NameAsIdentifier,

	// snowflake_stream_on_*: SchemaObjectIdentifier (all variants).
	"snowflake_stream_on_directory_table": SchemaObjectIdentifier(),
	"snowflake_stream_on_external_table":  SchemaObjectIdentifier(),
	"snowflake_stream_on_table":           SchemaObjectIdentifier(),
	"snowflake_stream_on_view":            SchemaObjectIdentifier(),

	// snowflake_streamlit: SchemaObjectIdentifier.
	"snowflake_streamlit": SchemaObjectIdentifier(),

	// snowflake_tag: SchemaObjectIdentifier.
	"snowflake_tag": SchemaObjectIdentifier(),

	// snowflake_tag_association: fixed 3-part compound ID
	// (<tag_id>|<tag_value>|<object_type>). All three parts are direct
	// TF parameters — no parameter reconstruction needed.
	// IdentifierFromProvider is sufficient because the full compound ID
	// is the resource ID on Create and the import string on Import.
	"snowflake_tag_association": config.IdentifierFromProvider,

	// snowflake_task: SchemaObjectIdentifier.
	"snowflake_task": SchemaObjectIdentifier(),

	// snowflake_user: bare name.
	"snowflake_user": config.NameAsIdentifier,

	// snowflake_user_programmatic_access_token: ID is '"<user>"|"<name>"' —
	// compound of user-supplied parameters (user + token name).
	"snowflake_user_programmatic_access_token": config.TemplatedStringAsIdentifier(
		"name",
		"\"{{ .parameters.user }}\"|\"{{ .external_name }}\"",
	),

	// snowflake_user_session_policy_attachment: ID is
	// '"<user>"|"<db>"."<schema>"."<session_policy>"' —
	// compound of user (external name) and session_policy_name field.
	"snowflake_user_session_policy_attachment": config.TemplatedStringAsIdentifier(
		"user_name",
		"\"{{ .external_name }}\"|{{ .parameters.session_policy_name }}",
	),

	// snowflake_view: SchemaObjectIdentifier.
	"snowflake_view": SchemaObjectIdentifier(),

	// snowflake_warehouse: bare name.
	"snowflake_warehouse": config.NameAsIdentifier,

	// =========================================================================
	// Preview resources (ShortGroup: preview)
	// =========================================================================

	// snowflake_account_authentication_policy_attachment: FQN stored in
	// authentication_policy field doubles as the identifier.
	"snowflake_account_authentication_policy_attachment": config.ParameterAsIdentifier("authentication_policy"),

	// snowflake_account_password_policy_attachment: FQN stored in
	// password_policy field doubles as the identifier.
	"snowflake_account_password_policy_attachment": config.ParameterAsIdentifier("password_policy"),

	// snowflake_alert: pipe-separated SchemaObjectIdentifier (EncodeSnowflakeID).
	"snowflake_alert": PipeSeparatedIdentifier(),

	// snowflake_api_integration: bare name.
	"snowflake_api_integration": config.NameAsIdentifier,

	// snowflake_api_integration_*: integration names (all variants).
	"snowflake_api_integration_amazon_api_gateway":          config.NameAsIdentifier,
	"snowflake_api_integration_azure_api_management":        config.NameAsIdentifier,
	"snowflake_api_integration_external_mcp_dynamic_client": config.NameAsIdentifier,
	"snowflake_api_integration_external_mcp_oauth2":         config.NameAsIdentifier,
	"snowflake_api_integration_git_repository_github_app":   config.NameAsIdentifier,
	"snowflake_api_integration_git_repository_oauth2":       config.NameAsIdentifier,
	"snowflake_api_integration_git_repository_private_link": config.NameAsIdentifier,
	"snowflake_api_integration_git_repository_token":        config.NameAsIdentifier,
	"snowflake_api_integration_google_cloud_api_gateway":    config.NameAsIdentifier,

	// snowflake_cortex_agent: SchemaObjectIdentifier (stable-style format).
	"snowflake_cortex_agent": SchemaObjectIdentifier(),

	// snowflake_cortex_search_service: pipe-separated SchemaObjectIdentifier.
	"snowflake_cortex_search_service": PipeSeparatedIdentifier(),

	// snowflake_dynamic_table: pipe-separated SchemaObjectIdentifier.
	"snowflake_dynamic_table": PipeSeparatedIdentifier(),

	// snowflake_email_notification_integration: bare name.
	"snowflake_email_notification_integration": config.NameAsIdentifier,

	// snowflake_external_function: ID includes argument signature
	// ("db"."schema"."func"(varchar,...)) — too complex for template.
	"snowflake_external_function": config.IdentifierFromProvider,

	// snowflake_external_table: pipe-separated SchemaObjectIdentifier.
	"snowflake_external_table": PipeSeparatedIdentifier(),

	// snowflake_failover_group: bare name.
	"snowflake_failover_group": config.NameAsIdentifier,

	// snowflake_file_format: pipe-separated SchemaObjectIdentifier.
	"snowflake_file_format": PipeSeparatedIdentifier(),

	// snowflake_function_*: ID includes name and argument types
	// ("db"."schema"."func"(varchar)) — signature from args block
	// makes template-based construction impractical.
	"snowflake_function_java":       config.IdentifierFromProvider,
	"snowflake_function_javascript": config.IdentifierFromProvider,
	"snowflake_function_python":     config.IdentifierFromProvider,
	"snowflake_function_scala":      config.IdentifierFromProvider,
	"snowflake_function_sql":        config.IdentifierFromProvider,

	// snowflake_iceberg_table_*: SchemaObjectIdentifier (stable-style format).
	"snowflake_iceberg_table_from_delta_files": SchemaObjectIdentifier(),
	"snowflake_iceberg_table_from_files":       SchemaObjectIdentifier(),

	// snowflake_job_service: SchemaObjectIdentifier (stable-style format).
	"snowflake_job_service": SchemaObjectIdentifier(),

	// snowflake_managed_account: bare name.
	"snowflake_managed_account": config.NameAsIdentifier,

	// snowflake_materialized_view: pipe-separated SchemaObjectIdentifier.
	"snowflake_materialized_view": PipeSeparatedIdentifier(),

	// snowflake_network_policy_attachment: ID is the network policy name
	// with '_attachment' suffix semantics — stored in network_policy_name.
	"snowflake_network_policy_attachment": config.ParameterAsIdentifier("network_policy_name"),

	// snowflake_notebook: SchemaObjectIdentifier (stable-style format).
	"snowflake_notebook": SchemaObjectIdentifier(),

	// snowflake_notification_integration: bare name.
	"snowflake_notification_integration": config.NameAsIdentifier,

	// snowflake_object_parameter: complex pipe-separated compound:
	// <key>|<object_type>|<object_identifier>.
	"snowflake_object_parameter": config.IdentifierFromProvider,

	// snowflake_pipe: pipe-separated SchemaObjectIdentifier.
	"snowflake_pipe": PipeSeparatedIdentifier(),

	// snowflake_postgres_instance: bare name.
	"snowflake_postgres_instance": config.NameAsIdentifier,

	// snowflake_procedure_*: ID includes name and argument types
	// ("db"."schema"."func"(varchar)) — same as function_* pattern.
	"snowflake_procedure_java":       config.IdentifierFromProvider,
	"snowflake_procedure_javascript": config.IdentifierFromProvider,
	"snowflake_procedure_python":     config.IdentifierFromProvider,
	"snowflake_procedure_scala":      config.IdentifierFromProvider,
	"snowflake_procedure_sql":        config.IdentifierFromProvider,

	// snowflake_semantic_view: SchemaObjectIdentifier (stable-style format).
	"snowflake_semantic_view": SchemaObjectIdentifier(),

	// snowflake_sequence: pipe-separated SchemaObjectIdentifier.
	"snowflake_sequence": PipeSeparatedIdentifier(),

	// snowflake_share: bare name.
	"snowflake_share": config.NameAsIdentifier,

	// snowflake_stage: pipe-separated SchemaObjectIdentifier.
	"snowflake_stage": PipeSeparatedIdentifier(),

	// snowflake_storage_integration: bare name.
	"snowflake_storage_integration": config.NameAsIdentifier,

	// snowflake_storage_lifecycle_policy: SchemaObjectIdentifier.
	"snowflake_storage_lifecycle_policy": SchemaObjectIdentifier(),

	// snowflake_table: pipe-separated SchemaObjectIdentifier.
	"snowflake_table": PipeSeparatedIdentifier(),

	// snowflake_table_column_masking_policy_application: complex compound of
	// table FQN, column name, and masking policy FQN.
	"snowflake_table_column_masking_policy_application": config.IdentifierFromProvider,

	// snowflake_table_constraint: uses '❄️' (snowflake) delimiter —
	// encodes constraint type, columns, and table.
	"snowflake_table_constraint": config.IdentifierFromProvider,

	// snowflake_table_storage_lifecycle_policy_attachment: two FQNs joined
	// by pipe — table FQN and storage lifecycle policy FQN.
	"snowflake_table_storage_lifecycle_policy_attachment": config.IdentifierFromProvider,

	// snowflake_user_authentication_policy_attachment: ID is
	// '"<user>"|"<db>"."<schema>"."<authentication_policy>"' —
	// compound of user (external name) and authentication_policy_name field.
	"snowflake_user_authentication_policy_attachment": config.TemplatedStringAsIdentifier(
		"user_name",
		"\"{{ .external_name }}\"|{{ .parameters.authentication_policy_name }}",
	),

	// snowflake_user_password_policy_attachment: ID is
	// '"<user>"|"<db>"."<schema>"."<password_policy>"' —
	// compound of user (external name) and password_policy_name field.
	"snowflake_user_password_policy_attachment": config.TemplatedStringAsIdentifier(
		"user_name",
		"\"{{ .external_name }}\"|{{ .parameters.password_policy_name }}",
	),

	// snowflake_user_public_keys: bare name (user name).
	"snowflake_user_public_keys": config.NameAsIdentifier,

	// snowflake_warehouse_adaptive: bare name.
	"snowflake_warehouse_adaptive": config.NameAsIdentifier,
}

// ExternalNameConfigurations applies all external name configs listed in the
// table ExternalNameConfigs and sets the version of those resources to v1beta1
// assuming they will be tested.
func ExternalNameConfigurations() config.ResourceOption {
	return func(r *config.Resource) {
		if e, ok := ExternalNameConfigs[r.Name]; ok {
			r.ExternalName = e
		}
	}
}

// ExternalNameConfigured returns the list of all resources whose external name
// is configured manually.
func ExternalNameConfigured() []string {
	l := make([]string, len(ExternalNameConfigs))
	i := 0
	for name := range ExternalNameConfigs {
		// $ is added to match the exact string since the format is regex.
		l[i] = name + "$"
		i++
	}
	return l
}
