package config

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/crossplane/upjet/v2/pkg/config"
)

// normalizeSFObjectID re-quotes a Snowflake object identifier. The input can be a
// bare name ("role"), an already-quoted name (`"role"`), or a fully-qualified name
// ("db"."role"). The output matches the form the TF provider's SDK emits from
// AccountObjectIdentifier/DatabaseObjectIdentifier.FullyQualifiedName(). Each
// dot-segment is wrapped in double quotes.
// The function splits on every ".". A quoted segment with a literal dot
// ("my.db"."role") splits wrongly. Snowflake object names with dots are rare.
// Use a quote-aware splitter if such a name shows up in practice.
func normalizeSFObjectID(raw string) string {
	segments := strings.Split(raw, ".")
	for i, seg := range segments {
		segments[i] = `"` + strings.Trim(seg, `"`) + `"`
	}
	return strings.Join(segments, ".")
}

// buildGrantAccountRoleID constructs the import ID for snowflake_grant_account_role
// from the resource parameters, matching helpers.EncodeSnowflakeID in the TF
// provider (pkg/resources/grant_account_role.go): '<role>|ROLE|<parent_role>'
// or '<role>|USER|<user>'.
func buildGrantAccountRoleID(parameters map[string]any) (string, error) {
	roleName, _ := parameters["role_name"].(string)
	if roleName == "" {
		return "", fmt.Errorf("grant_account_role: role_name is required")
	}
	var objectType, target string
	if v, _ := parameters["parent_role_name"].(string); v != "" {
		objectType, target = "ROLE", v
	} else if v, _ := parameters["user_name"].(string); v != "" {
		objectType, target = "USER", v
	} else {
		return "", fmt.Errorf("grant_account_role: neither parent_role_name nor user_name is set")
	}
	return strings.Join([]string{normalizeSFObjectID(roleName), objectType, normalizeSFObjectID(target)}, "|"), nil
}

// GrantAccountRoleIdentifier returns an ExternalName for snowflake_grant_account_role.
// No field is omitted from spec.forProvider. role_name, parent_role_name, and
// user_name all stay visible, so multiple grants can share a role_name. This
// matches the TF for_each pattern because the K8s identity does not depend on
// any single field.
func GrantAccountRoleIdentifier() config.ExternalName {
	return config.NewExternalNameFrom(config.IdentifierFromProvider,
		config.WithGetIDFn(func(fn config.GetIDFn, ctx context.Context, externalName string, parameters map[string]any, providerConfig map[string]any) (string, error) {
			if id, err := buildGrantAccountRoleID(parameters); err == nil && id != "" {
				return id, nil
			}
			return fn(ctx, externalName, parameters, providerConfig)
		}),
	)
}

// buildGrantApplicationRoleID constructs the import ID for
// snowflake_grant_application_role, matching helpers.EncodeResourceIdentifier
// in the TF provider (pkg/resources/grant_application_role.go):
// '<app_role_fqn>|ACCOUNT_ROLE|<parent_account_role>' or
// '<app_role_fqn>|APPLICATION|<application>'.
func buildGrantApplicationRoleID(parameters map[string]any) (string, error) {
	appRoleName, _ := parameters["application_role_name"].(string)
	if appRoleName == "" {
		return "", fmt.Errorf("grant_application_role: application_role_name is required")
	}
	var objectType, target string
	if v, _ := parameters["parent_account_role_name"].(string); v != "" {
		objectType, target = "ACCOUNT_ROLE", v
	} else if v, _ := parameters["application_name"].(string); v != "" {
		objectType, target = "APPLICATION", v
	} else {
		return "", fmt.Errorf("grant_application_role: neither parent_account_role_name nor application_name is set")
	}
	return strings.Join([]string{normalizeSFObjectID(appRoleName), objectType, normalizeSFObjectID(target)}, "|"), nil
}

// GrantApplicationRoleIdentifier returns an ExternalName for
// snowflake_grant_application_role. See GrantAccountRoleIdentifier for why no
// field is omitted.
func GrantApplicationRoleIdentifier() config.ExternalName {
	return config.NewExternalNameFrom(config.IdentifierFromProvider,
		config.WithGetIDFn(func(fn config.GetIDFn, ctx context.Context, externalName string, parameters map[string]any, providerConfig map[string]any) (string, error) {
			if id, err := buildGrantApplicationRoleID(parameters); err == nil && id != "" {
				return id, nil
			}
			return fn(ctx, externalName, parameters, providerConfig)
		}),
	)
}

// buildGrantDatabaseRoleID constructs the import ID for
// snowflake_grant_database_role, matching helpers.EncodeResourceIdentifier in
// the TF provider (pkg/resources/grant_database_role.go):
// '<db_role_fqn>|ROLE|<parent_role>' or
// '<db_role_fqn>|DATABASE ROLE|<parent_db_role>' or
// '<db_role_fqn>|SHARE|<share>'.
// Note: sdk.ObjectTypeDatabaseRole.String() is "DATABASE ROLE" (with a space),
// not "DATABASE_ROLE". Verified against pkg/sdk/object_types.go.
func buildGrantDatabaseRoleID(parameters map[string]any) (string, error) {
	dbRoleName, _ := parameters["database_role_name"].(string)
	if dbRoleName == "" {
		return "", fmt.Errorf("grant_database_role: database_role_name is required")
	}
	var objectType, target string
	if v, _ := parameters["parent_role_name"].(string); v != "" {
		objectType, target = "ROLE", v
	} else if v, _ := parameters["parent_database_role_name"].(string); v != "" {
		objectType, target = "DATABASE ROLE", v
	} else if v, _ := parameters["share_name"].(string); v != "" {
		objectType, target = "SHARE", v
	} else {
		return "", fmt.Errorf("grant_database_role: none of parent_role_name, parent_database_role_name, share_name is set")
	}
	return strings.Join([]string{normalizeSFObjectID(dbRoleName), objectType, normalizeSFObjectID(target)}, "|"), nil
}

// GrantDatabaseRoleIdentifier returns an ExternalName for
// snowflake_grant_database_role. See GrantAccountRoleIdentifier for why no
// field is omitted.
func GrantDatabaseRoleIdentifier() config.ExternalName {
	return config.NewExternalNameFrom(config.IdentifierFromProvider,
		config.WithGetIDFn(func(fn config.GetIDFn, ctx context.Context, externalName string, parameters map[string]any, providerConfig map[string]any) (string, error) {
			if id, err := buildGrantDatabaseRoleID(parameters); err == nil && id != "" {
				return id, nil
			}
			return fn(ctx, externalName, parameters, providerConfig)
		}),
	)
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
// The single-object case ALWAYS carries the OnObject sub-type marker:
// ParseGrantPrivilegesTo{Account,Database}RoleId only accepts the 8-part
// "…|OnSchemaObject|OnObject|<object_type>|<object_name>" form.
// grant_privileges_to_database_role_identifier.go errors
// "invalid OnSchemaObjectGrantKind" for any other form. The upstream encoder
// (grant_privileges_identifier_commons.go OnSchemaObjectGrantData.String)
// emits OnObject unconditionally. all_privileges does not change the ID.
func onSchemaObjectBlockSuffix(so map[string]any) (string, []string, bool, error) {
	// Single object: object_type + object_name
	if objType, _ := so["object_type"].(string); objType != "" {
		objName, _ := so["object_name"].(string)
		return "OnObject", []string{objType, objName}, true, nil
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

const strTrue = "true"

// paramBool reads a Terraform parameter that may be bool or "true"/"false" string.
func paramBool(parameters map[string]any, key string) bool {
	v, ok := parameters[key]
	if !ok {
		return false
	}
	switch b := v.(type) {
	case bool:
		return b
	case string:
		return b == strTrue
	default:
		return false
	}
}

// grantPrivilegesBaseStr constructs the common base prefix for grant_privileges
// resources: <role_name>|<with_grant_option>|<always_apply>|<privileges>.
// privileges is formatted as comma-separated sorted list, or "ALL" when
// all_privileges is true.
func grantPrivilegesBaseStr(roleName string, parameters map[string]any) string {
	wgo := "false"
	if paramBool(parameters, "with_grant_option") {
		wgo = strTrue
	}
	aa := "false"
	if paramBool(parameters, "always_apply") {
		aa = strTrue
	}
	var privs string
	if paramBool(parameters, "all_privileges") {
		privs = "ALL"
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

// grantOwnershipOnBlockSuffix processes the "on" block for snowflake_grant_ownership
// and returns the grant type and suffix parts.
func grantOwnershipOnBlockSuffix(onBlock map[string]any) (string, []string, error) {
	// OnObject: object_type + object_name
	if objType, _ := onBlock["object_type"].(string); objType != "" {
		objName, _ := onBlock["object_name"].(string)
		return "OnObject", []string{objType, objName}, nil
	}
	// OnAll: all[0] sub-block (object_type_plural + in_database or in_schema)
	if allRaw, _ := onBlock["all"].([]any); len(allRaw) > 0 {
		allBlock, ok := allRaw[0].(map[string]any)
		if !ok {
			return "", nil, fmt.Errorf("grant_ownership: invalid 'all' block")
		}
		parts, err := grantAllFutureParts(allBlock)
		if err != nil {
			return "", nil, fmt.Errorf("grant_ownership: 'all' block: %w", err)
		}
		return "OnAll", parts, nil
	}
	// OnFuture: future[0] sub-block (object_type_plural + in_database or in_schema)
	if futureRaw, _ := onBlock["future"].([]any); len(futureRaw) > 0 {
		futureBlock, ok := futureRaw[0].(map[string]any)
		if !ok {
			return "", nil, fmt.Errorf("grant_ownership: invalid 'future' block")
		}
		parts, err := grantAllFutureParts(futureBlock)
		if err != nil {
			return "", nil, fmt.Errorf("grant_ownership: 'future' block: %w", err)
		}
		return "OnFuture", parts, nil
	}
	return "", nil, fmt.Errorf("grant_ownership: unable to determine grant type from 'on' block")
}

// GrantOwnershipIdentifier returns an ExternalName for snowflake_grant_ownership.
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

	onRaw, _ := parameters["on"].([]any)
	if len(onRaw) == 0 {
		return "", fmt.Errorf("grant_ownership: 'on' block is required")
	}
	onBlock, ok := onRaw[0].(map[string]any)
	if !ok {
		return "", fmt.Errorf("grant_ownership: invalid 'on' block")
	}

	grantType, suffix, err := grantOwnershipOnBlockSuffix(onBlock)
	if err != nil {
		return "", err
	}
	parts := append([]string{roleType, roleID, outboundPrivileges, grantType}, suffix...)
	return strings.Join(parts, "|"), nil
}

// =============================================================================
// Snowflake-specific external name identifier functions
// processOnSchemaID handles the on_schema block for grant_privileges resources.
func processOnSchemaID(parameters map[string]any, base, paramName, resourceName string) (string, bool, error) {
	schemaRaw, _ := parameters[paramName].([]any)
	if len(schemaRaw) == 0 {
		return "", false, nil
	}
	schemaBlock, ok := schemaRaw[0].(map[string]any)
	if !ok {
		return "", false, fmt.Errorf("%s: invalid '%s' block", resourceName, paramName)
	}
	subType, suffix, found := onSchemaBlockSuffix(schemaBlock)
	if !found {
		return "", false, fmt.Errorf("%s: %s block variant not recognized", resourceName, paramName)
	}
	return fmt.Sprintf("%s|OnSchema|%s|%s", base, subType, strings.Join(suffix, "|")), true, nil
}

// processOnSchemaObjectID handles the on_schema_object block for grant_privileges resources.
func processOnSchemaObjectID(parameters map[string]any, base, resourceName string) (string, bool, error) {
	soRaw, _ := parameters["on_schema_object"].([]any)
	if len(soRaw) == 0 {
		return "", false, nil
	}
	soBlock, ok := soRaw[0].(map[string]any)
	if !ok {
		return "", false, fmt.Errorf("%s: invalid 'on_schema_object' block", resourceName)
	}
	subType, suffix, found, err := onSchemaObjectBlockSuffix(soBlock)
	if err != nil {
		return "", false, fmt.Errorf("%s: %w", resourceName, err)
	}
	if !found {
		return "", false, fmt.Errorf("%s: on_schema_object block variant not recognized", resourceName)
	}
	return fmt.Sprintf("%s|OnSchemaObject|%s|%s", base, subType, strings.Join(suffix, "|")), true, nil
}

// GrantPrivilegesToAccountRoleIdentifier returns an ExternalName for
// snowflake_grant_privileges_to_account_role.
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

	base := grantPrivilegesBaseStr(roleName, parameters)

	if paramBool(parameters, "on_account") {
		return base + "|OnAccount", nil
	}

	if oaRaw, _ := parameters["on_account_object"].([]any); len(oaRaw) > 0 {
		oaBlock, ok := oaRaw[0].(map[string]any)
		if !ok {
			return "", fmt.Errorf("grant_privileges_to_account_role: invalid 'on_account_object' block")
		}
		objType, _ := oaBlock["object_type"].(string)
		objName, _ := oaBlock["object_name"].(string)
		return fmt.Sprintf("%s|OnAccountObject|%s|%s", base, objType, objName), nil
	}

	if id, found, err := processOnSchemaID(parameters, base, "on_schema", "grant_privileges_to_account_role"); err != nil {
		return "", err
	} else if found {
		return id, nil
	}

	if id, found, err := processOnSchemaObjectID(parameters, base, "grant_privileges_to_account_role"); err != nil {
		return "", err
	} else if found {
		return id, nil
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

	base := grantPrivilegesBaseStr(roleName, parameters)

	// OnDatabase: string field
	if db, _ := parameters["on_database"].(string); db != "" {
		return fmt.Sprintf("%s|OnDatabase|%s", base, db), nil
	}

	if id, found, err := processOnSchemaID(parameters, base, "on_schema", "grant_privileges_to_database_role"); err != nil {
		return "", err
	} else if found {
		return id, nil
	}

	if id, found, err := processOnSchemaObjectID(parameters, base, "grant_privileges_to_database_role"); err != nil {
		return "", err
	} else if found {
		return id, nil
	}

	return "", fmt.Errorf("grant_privileges_to_database_role: one of on_database, on_schema, or on_schema_object is required")
}

// GrantPrivilegesToShareIdentifier returns an ExternalName for
// snowflake_grant_privileges_to_share.
//
// The TF resource ID is a compound, variable-length pipe-separated string:
//
//	<to_share>|<privileges>|<grant_type>|<grant_identifier>
//
// No field is omitted from spec.forProvider. to_share stays as a regular
// spec parameter, so multiple grants under the same share are possible
// with different on_* targets. This matches the TF for_each pattern.
func GrantPrivilegesToShareIdentifier() config.ExternalName {
	return config.NewExternalNameFrom(
		config.IdentifierFromProvider,
		config.WithGetIDFn(func(fn config.GetIDFn, ctx context.Context, externalName string, parameters map[string]any, providerConfig map[string]any) (string, error) {
			if id, err := buildGrantPrivilegesToShareID(parameters); err == nil && id != "" {
				return id, nil
			}
			return fn(ctx, externalName, parameters, providerConfig)
		}),
	)
}

func formatSharePrivileges(parameters map[string]any) string {
	privRaw, ok := parameters["privileges"].([]any)
	if !ok || len(privRaw) == 0 {
		return ""
	}
	p := make([]string, len(privRaw))
	for i, v := range privRaw {
		p[i] = fmt.Sprint(v)
	}
	slices.Sort(p)
	return strings.Join(p, ",")
}

// shareOnFields lists the ExactlyOneOf on_* fields for
// snowflake_grant_privileges_to_share, in the order the TF provider checks
// them, paired with their ID label.
var shareOnFields = []struct{ key, label string }{
	{"on_database", "OnDatabase"},
	{"on_function", "OnFunction"},
	{"on_schema", "OnSchema"},
	{"on_table", "OnTable"},
	{"on_all_tables_in_schema", "OnAllTablesInSchema"},
	{"on_tag", "OnTag"},
	{"on_view", "OnView"},
}

func buildGrantPrivilegesToShareID(parameters map[string]any) (string, error) {
	shareName, _ := parameters["to_share"].(string)
	if shareName == "" {
		return "", fmt.Errorf("grant_privileges_to_share: to_share is required")
	}

	// Inline privileges formatting, not a call to grantPrivilegesBaseStr.
	// Share grants have no with_grant_option or always_apply fields.
	privs := formatSharePrivileges(parameters)
	if privs == "" {
		return "", fmt.Errorf("grant_privileges_to_share: privileges is required")
	}

	base := fmt.Sprintf("%s|%s", shareName, privs)

	// Values come from spec.forProvider as the user provides them;
	// the TF provider's SDK handles both bare and quoted forms.
	for _, f := range shareOnFields {
		if v, _ := parameters[f.key].(string); v != "" {
			return fmt.Sprintf("%s|%s|%s", base, f.label, v), nil
		}
	}

	return "", fmt.Errorf("grant_privileges_to_share: one of on_database, on_function, on_schema, on_table, on_all_tables_in_schema, on_tag, or on_view is required")
}
