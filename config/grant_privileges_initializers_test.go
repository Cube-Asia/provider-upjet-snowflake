package config

import (
	"context"
	"testing"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestParseGrantPrivilegesBaseID(t *testing.T) {
	cases := map[string]struct {
		id                                              string
		wantRole                                        string
		wantWGO, wantAlwaysApply, wantAllPrivileges, ok bool
		wantPrivileges                                  []string
	}{
		"on account, explicit privileges": {
			id:             "parent-minimal-role|false|false|CREATE DATABASE,CREATE USER|OnAccount",
			wantRole:       "parent-minimal-role",
			ok:             true,
			wantPrivileges: []string{"CREATE DATABASE", "CREATE USER"},
		},
		"with grant option and always apply": {
			id:                "r|true|true|ALL|OnAccount",
			wantRole:          "r",
			wantWGO:           true,
			wantAlwaysApply:   true,
			wantAllPrivileges: true,
			ok:                true,
		},
		"malformed": {
			id: "not-enough-parts",
			ok: false,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			role, wgo, aa, all, privs, ok := parseGrantPrivilegesBaseID(tc.id)
			if ok != tc.ok {
				t.Fatalf("ok = %v, want %v", ok, tc.ok)
			}
			if !ok {
				return
			}
			if role != tc.wantRole || wgo != tc.wantWGO || aa != tc.wantAlwaysApply || all != tc.wantAllPrivileges {
				t.Fatalf("got (%q,%v,%v,%v), want (%q,%v,%v,%v)", role, wgo, aa, all, tc.wantRole, tc.wantWGO, tc.wantAlwaysApply, tc.wantAllPrivileges)
			}
			if len(privs) != len(tc.wantPrivileges) {
				t.Fatalf("privileges = %v, want %v", privs, tc.wantPrivileges)
			}
			for i := range privs {
				if privs[i] != tc.wantPrivileges[i] {
					t.Fatalf("privileges = %v, want %v", privs, tc.wantPrivileges)
				}
			}
		})
	}
}

// testGrantPrivilegesSchema mirrors the fields of
// grantPrivilegesToAccountRoleSchema that backfillGrantPrivilegesFields
// touches, enough to exercise d.Set/d.Get through *schema.ResourceData.
func testGrantPrivilegesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"account_role_name":           {Type: schema.TypeString, Optional: true},
		"with_grant_option":           {Type: schema.TypeBool, Optional: true},
		"always_apply":                {Type: schema.TypeBool, Optional: true},
		"all_privileges":              {Type: schema.TypeBool, Optional: true},
		"strict_privilege_management": {Type: schema.TypeBool, Optional: true},
		"privileges": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
}

// TestWrapGrantPrivilegesReadContext_BackfillsMissingForceNewField locks the
// actual bug: with_grant_option is ForceNew and the real Read() never sets
// it, only Import does. The wrapped ReadContext must backfill it (and its
// siblings) from the resource's own ID after every real Read, regardless of
// whether the incoming ResourceData already had it set.
func TestWrapGrantPrivilegesReadContext_BackfillsMissingForceNewField(t *testing.T) {
	sc := testGrantPrivilegesSchema()
	var realReadCalls int
	realRead := func(_ context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
		realReadCalls++
		// Mirrors the real upstream Read(): sets privileges, never touches
		// with_grant_option/all_privileges/always_apply/strict_privilege_management.
		return diag.FromErr(d.Set("privileges", []string{"CREATE DATABASE", "CREATE USER"}))
	}

	r := &schema.Resource{Schema: sc, ReadContext: realRead}
	wrapGrantPrivilegesReadContext("account_role_name")(&ujconfig.Resource{TerraformResource: r})

	d := schema.TestResourceDataRaw(t, sc, map[string]any{})
	d.SetId("parent-minimal-role|false|false|CREATE DATABASE,CREATE USER|OnAccount")

	diags := r.ReadContext(context.Background(), d, nil)
	if diags.HasError() {
		t.Fatalf("ReadContext() diags = %v", diags)
	}
	if realReadCalls != 1 {
		t.Fatalf("real Read called %d times, want 1", realReadCalls)
	}
	if got := d.Get("with_grant_option").(bool); got != false {
		t.Fatalf("with_grant_option = %v, want false", got)
	}
	if got := d.Get("account_role_name").(string); got != "parent-minimal-role" {
		t.Fatalf("account_role_name = %q, want %q", got, "parent-minimal-role")
	}
	if got := d.Get("privileges").(*schema.Set).Len(); got != 2 {
		t.Fatalf("privileges len = %d, want 2 (real Read's value must survive)", got)
	}

	// A resource found deleted upstream (Read clears the ID) must not panic
	// or attempt to parse an empty ID.
	deletedRead := func(_ context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
		d.SetId("")
		return nil
	}
	r2 := &schema.Resource{Schema: sc, ReadContext: deletedRead}
	wrapGrantPrivilegesReadContext("account_role_name")(&ujconfig.Resource{TerraformResource: r2})
	d2 := schema.TestResourceDataRaw(t, sc, map[string]any{})
	d2.SetId("parent-minimal-role|false|false|CREATE DATABASE,CREATE USER|OnAccount")
	if diags := r2.ReadContext(context.Background(), d2, nil); diags.HasError() {
		t.Fatalf("ReadContext() diags = %v", diags)
	}
	if d2.Id() != "" {
		t.Fatalf("Id() = %q, want empty after upstream deletion", d2.Id())
	}
}
