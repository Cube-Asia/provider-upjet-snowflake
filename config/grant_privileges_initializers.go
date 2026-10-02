package config

import (
	"context"
	"strings"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ---------------------------------------------------------------------------
// snowflake_grant_privileges_to_{account,database}_role Read() gap workaround
//
// The Read function of the upstream Snowflake TF provider sets only
// d.Set("privileges", ...) (and conditionally always_apply_trigger) on every
// refresh. It never sets with_grant_option, all_privileges, always_apply, or
// strict_privilege_management again. Only Import sets them, and Import reads
// them from the compound ID (see
// .work/snowflakedb/snowflake/pkg/resources/grant_privileges_to_{account,database}_role.go,
// Read* vs Import* functions). This backfill is not a live-drift check. Read
// does not query Snowflake again for the grant-option state. Even Import
// derives these fields from the ID string, never from a live API call. So a
// backfill from the ID on every Read matches what Import already does. The
// only difference is that it runs on every refresh instead of once.
//
// Why a Read wrapper, and not something further up the stack:
// upjet's no-fork Observe() calls schema.Resource.RefreshWithoutUpgrade.
// That call uses an in-memory per-resource cache (upjet's
// OperationTrackerStore). Create/Update/Observe populate the cache once. The
// provider process then reuses it for its whole life. atProvider (status) is
// consulted only to rebuild the cache when it is cold (for example right
// after a provider restart). A managed.Initializer that backfills
// status.atProvider therefore helps only right after a restart. A resource
// created and observed within the same long-running process keeps seeing the
// incomplete cached state forever. Observe computes its diff against the
// FRESHLY REFRESHED state (diffState = newState, set right after
// RefreshWithoutUpgrade, before the diff call). So wrapping ReadContext
// fixes the state before every diff, no matter which path (cold
// reconstruction or warm cache) fed it in.
//
// with_grant_option is ForceNew in the TF schema. If the field stays unset,
// upjet's assertNoForceNew guard fails forever with: "refuse to update the
// external resource ... requires replacing it: cannot change the value of
// the argument with_grant_option". SYNCED never recovers on its own.
//
// The wrap goes through *schema.Resource.ReadContext. The SDK supports this
// mutable function field on the object github.com/Snowflake-Labs/... hands
// us via ujconfig.WithTerraformProvider. No vendored .go file is edited.
// ---------------------------------------------------------------------------

// wrapGrantPrivilegesReadContext returns a config.ResourceOption that wraps
// r.TerraformResource.ReadContext, backfilling the fields the real Read never
// sets. roleNameKey is "account_role_name" or "database_role_name".
func wrapGrantPrivilegesReadContext(roleNameKey string) func(r *ujconfig.Resource) {
	return func(r *ujconfig.Resource) {
		orig := r.TerraformResource.ReadContext
		if orig == nil {
			return
		}
		r.TerraformResource.ReadContext = func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
			diags := orig(ctx, d, meta)
			if diags.HasError() || d.Id() == "" {
				// This is a real error, or Read found the resource deleted
				// upstream and cleared the ID. There is nothing to backfill.
				return diags
			}
			backfillGrantPrivilegesFields(d, roleNameKey)
			return diags
		}
	}
}

// backfillGrantPrivilegesFields sets the fields Read() forgets. It decodes
// them from the resource's own compound ID (see parseGrantPrivilegesBaseID).
// privileges is deliberately left untouched. Read() already recomputes it
// correctly from a live "SHOW GRANTS" call.
//
// strict_privilege_management is pinned to false on every refresh, even when
// the spec asks for true. The pin covers account-role grants only. This
// mirrors upstream Import exactly (grant_privileges_to_account_role.go:
// d.Set("strict_privilege_management", false). It is not encoded in the
// compound ID, so Import cannot recover it. The database-role resource has
// no such attribute in the vendored v2.19.0 schema. d.Set on a missing key
// makes the SDK log "[ERROR] setting state: Invalid address to set" on every
// Read. So the pin is gated on roleNameKey. If a provider bump adds the
// field to the database-role schema, mirror that version's Import behavior
// here in the same change. Setting it to true additionally requires the
// provider-level GRANTS_STRICT_PRIVILEGE_MANAGEMENT experimental feature. No
// known claim enables it. If a claim ever needs it, the backfill must learn
// to respect a set state value instead of mirroring Import.
func backfillGrantPrivilegesFields(d *schema.ResourceData, roleNameKey string) {
	roleName, withGrantOption, alwaysApply, allPrivileges, _, ok := parseGrantPrivilegesBaseID(d.Id())
	if !ok {
		return
	}
	_ = d.Set(roleNameKey, roleName)
	_ = d.Set("with_grant_option", withGrantOption)
	_ = d.Set("always_apply", alwaysApply)
	_ = d.Set("all_privileges", allPrivileges)
	// account_role_name only: see the strict_privilege_management note above.
	if roleNameKey == accountRoleNameKey {
		_ = d.Set("strict_privilege_management", false)
	}
}

// parseGrantPrivilegesBaseID decodes the shared
// <role>|<with_grant_option>|<always_apply>|<privileges>|... prefix used by
// both snowflake_grant_privileges_to_account_role and
// snowflake_grant_privileges_to_database_role compound IDs. See
// grantPrivilegesBaseStr in grant_id_builders.go for the encoder.
func parseGrantPrivilegesBaseID(id string) (roleName string, withGrantOption, alwaysApply, allPrivileges bool, privileges []string, ok bool) {
	parts := strings.SplitN(id, "|", 5)
	if len(parts) < 4 || parts[0] == "" {
		return "", false, false, false, nil, false
	}
	roleName = parts[0]
	withGrantOption = parts[1] == strTrue
	alwaysApply = parts[2] == strTrue
	switch parts[3] {
	case "":
	case "ALL":
		allPrivileges = true
	default:
		privileges = strings.Split(parts[3], ",")
	}
	return roleName, withGrantOption, alwaysApply, allPrivileges, privileges, true
}

// GrantPrivilegesReadGapWorkaround wires the ReadContext wrapper above onto
// the two affected resources. It runs as a default resource option, so it
// covers the cluster-scoped and namespaced-scoped providers the same way.
func GrantPrivilegesReadGapWorkaround() ujconfig.ResourceOption {
	return func(r *ujconfig.Resource) {
		switch r.Name {
		case "snowflake_grant_privileges_to_account_role":
			wrapGrantPrivilegesReadContext(accountRoleNameKey)(r)
		case "snowflake_grant_privileges_to_database_role":
			wrapGrantPrivilegesReadContext(databaseRoleNameKey)(r)
		}
	}
}
