package resourcelist

import (
	"strings"
)

// ToKind converts a Snowflake Terraform resource name to a Crossplane Kind.
// It strips the "snowflake_" prefix and PascalCases the remainder.
//
// Examples:
//
//	snowflake_account_role  → AccountRole
//	snowflake_grant_privileges_to_account_role → GrantPrivilegesToAccountRole
//
// This is needed because upjet's default Kind derivation treats
// word[1] as a "group" suffix and drops it, causing collisions:
// snowflake_account_role and snowflake_database_role would both
// produce Kind "Role" without this fix.
func ToKind(tfResource string) string {
	name := strings.TrimPrefix(tfResource, "snowflake_")
	parts := strings.Split(name, "_")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, "")
}
