package config

import (
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
// Only appropriate when nameField genuinely IS the resource's sole K8s
// identity. For resources where nameField can repeat across many sibling
// grants — grant_account_role, grant_application_role, grant_database_role,
// and grant_privileges_to_share, which can grant a share multiple
// non-overlapping privilege sets across different on_* targets under one
// to_share — use a custom NewExternalNameFrom(IdentifierFromProvider, ...)
// builder instead, see GrantAccountRoleIdentifier, so the field stays a
// regular, non-omitted spec parameter instead of being forced into the
// K8s identity.
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
