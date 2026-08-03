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
// The upstream Snowflake TF provider's Read function for these two resources
// only calls d.Set("privileges", ...) (and conditionally
// always_apply_trigger) on every refresh. It never re-sets with_grant_option,
// all_privileges, always_apply, or strict_privilege_management — only
// Import does, from the compound ID (see
// .work/snowflakedb/snowflake/pkg/resources/grant_privileges_to_{account,database}_role.go,
// Read* vs Import* functions). This is not a live-drift check either way:
// Snowflake's grant-option state isn't re-queried by Read; even Import only
// derives these fields from the ID string, never from a live API call. So
// backfilling them from the ID on every Read is exactly what Import already
// does, just on every refresh instead of once.
//
// Why this has to be a Read wrapper and not something further up the stack:
// upjet's no-fork Observe() calls schema.Resource.RefreshWithoutUpgrade,
// seeded from an in-memory per-resource cache (upjet's OperationTrackerStore)
// that is populated once by Create/Update/Observe and then reused for the
// life of the provider process — atProvider (status) is only consulted to
// reconstruct that cache when it's cold (e.g. right after a provider
// restart). A managed.Initializer backfilling status.atProvider therefore
// only helps immediately after a restart; a resource created and observed
// within the same long-running process keeps hitting the incomplete cached
// state forever. Observe's diff is computed against the FRESHLY REFRESHED
// state (diffState = newState, set right after RefreshWithoutUpgrade, before
// the diff call) — so wrapping ReadContext fixes the state before every
// single diff, regardless of which path (cold reconstruction or warm cache)
// fed it in.
//
// with_grant_option is ForceNew in the TF schema, so leaving it unset
// permanently trips upjet's assertNoForceNew guard: "refuse to update the
// external resource ... requires replacing it: cannot change the value of
// the argument with_grant_option". SYNCED never recovers on its own.
//
// The wrap wires through *schema.Resource.ReadContext — an SDK-supported,
// mutable function field on the object github.com/Snowflake-Labs/... hands
// us via ujconfig.WithTerraformProvider — not an edit to any vendored .go
// file.
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
				// Real error, or the resource was found deleted upstream
				// (Read cleared the ID) — nothing to backfill.
				return diags
			}
			backfillGrantPrivilegesFields(d, roleNameKey)
			return diags
		}
	}
}

// backfillGrantPrivilegesFields sets the fields Read() forgets, decoded from
// the resource's own compound ID (see parseGrantPrivilegesBaseID). privileges
// is deliberately left untouched — Read() already recomputes it correctly
// from a live "SHOW GRANTS" call.
func backfillGrantPrivilegesFields(d *schema.ResourceData, roleNameKey string) {
	roleName, withGrantOption, alwaysApply, allPrivileges, _, ok := parseGrantPrivilegesBaseID(d.Id())
	if !ok {
		return
	}
	_ = d.Set(roleNameKey, roleName)
	_ = d.Set("with_grant_option", withGrantOption)
	_ = d.Set("always_apply", alwaysApply)
	_ = d.Set("all_privileges", allPrivileges)
	_ = d.Set("strict_privilege_management", false)
}

// parseGrantPrivilegesBaseID decodes the shared
// <role>|<with_grant_option>|<always_apply>|<privileges>|... prefix used by
// both snowflake_grant_privileges_to_account_role and
// snowflake_grant_privileges_to_database_role compound IDs. See
// grantPrivilegesBaseStr in external_name.go for the encoder.
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
// the two affected resources. Applied as a default resource option so it
// covers both the cluster-scoped and namespaced-scoped providers uniformly.
func GrantPrivilegesReadGapWorkaround() ujconfig.ResourceOption {
	return func(r *ujconfig.Resource) {
		switch r.Name {
		case "snowflake_grant_privileges_to_account_role":
			wrapGrantPrivilegesReadContext("account_role_name")(r)
		case "snowflake_grant_privileges_to_database_role":
			wrapGrantPrivilegesReadContext("database_role_name")(r)
		}
	}
}
