package config

import (
	"bytes"
	"context"
	"log"
	"os"
	"reflect"
	"strings"
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
		"empty role": {
			id: "|true|false|USAGE|OnDatabase|db",
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

// testGrantPrivilegesSchema copies the fields of
// grantPrivilegesTo{Account,Database}RoleSchema that
// backfillGrantPrivilegesFields touches. The copy is enough to exercise
// d.Set/d.Get through *schema.ResourceData. withStrictPrivilegeManagement
// reflects the vendored v2.19.0 schemas: the attribute exists on the
// account-role resource only.
func testGrantPrivilegesSchema(roleKey string, withStrictPrivilegeManagement bool) map[string]*schema.Schema {
	sc := map[string]*schema.Schema{
		roleKey:             {Type: schema.TypeString, Optional: true},
		"with_grant_option": {Type: schema.TypeBool, Optional: true},
		"always_apply":      {Type: schema.TypeBool, Optional: true},
		"all_privileges":    {Type: schema.TypeBool, Optional: true},
		"privileges": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
	if withStrictPrivilegeManagement {
		sc["strict_privilege_management"] = &schema.Schema{Type: schema.TypeBool, Optional: true}
	}
	return sc
}

// TestWrapGrantPrivilegesReadContext_BackfillsMissingForceNewField locks the
// actual bug. with_grant_option is ForceNew. The real Read() never sets it.
// Only Import sets it. The wrapped ReadContext must backfill it and its
// siblings from the resource's own ID after every real Read. It must do this
// even when the incoming ResourceData already has a value.
func TestWrapGrantPrivilegesReadContext_BackfillsMissingForceNewField(t *testing.T) {
	sc := testGrantPrivilegesSchema(accountRoleNameKey, true)
	var realReadCalls int
	realRead := func(_ context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
		realReadCalls++
		// Mirrors the real upstream Read(). It sets privileges and never
		// touches with_grant_option, all_privileges, always_apply, or
		// strict_privilege_management.
		return diag.FromErr(d.Set("privileges", []string{"CREATE DATABASE", "CREATE USER"}))
	}

	r := &schema.Resource{Schema: sc, ReadContext: realRead}
	wrapGrantPrivilegesReadContext(accountRoleNameKey)(&ujconfig.Resource{TerraformResource: r})

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
	if got := d.Get(accountRoleNameKey).(string); got != "parent-minimal-role" {
		t.Fatalf("account_role_name = %q, want %q", got, "parent-minimal-role")
	}
	if got := d.Get("privileges").(*schema.Set).Len(); got != 2 {
		t.Fatalf("privileges len = %d, want 2 (real Read's value must survive)", got)
	}

	// If the resource is gone upstream, Read clears the ID. The wrapper must
	// not panic and must not try to parse an empty ID.
	deletedRead := func(_ context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
		d.SetId("")
		return nil
	}
	r2 := &schema.Resource{Schema: sc, ReadContext: deletedRead}
	wrapGrantPrivilegesReadContext(accountRoleNameKey)(&ujconfig.Resource{TerraformResource: r2})
	d2 := schema.TestResourceDataRaw(t, sc, map[string]any{})
	d2.SetId("parent-minimal-role|false|false|CREATE DATABASE,CREATE USER|OnAccount")
	if diags := r2.ReadContext(context.Background(), d2, nil); diags.HasError() {
		t.Fatalf("ReadContext() diags = %v", diags)
	}
	if d2.Id() != "" {
		t.Fatalf("Id() = %q, want empty after upstream deletion", d2.Id())
	}
}

// The builders emit the external ID bare. IDs observed in a live Snowflake
// account arrive quoted, because Snowflake's ID delegation quotes each name
// segment. The parser must accept quoted segments verbatim. Grants created
// by this provider then keep parsing on the Read path after the ID
// round-trips through Snowflake.
func TestParseGrantPrivilegesBaseIDLiveWireFormat(t *testing.T) {
	const liveID = `"test-e2e-role"|true|false|USAGE|OnAccountObject|DATABASE|"test-e2e-db"`
	role, wgo, aa, all, privs, ok := parseGrantPrivilegesBaseID(liveID)
	if !ok {
		t.Fatalf("parse(%s) failed", liveID)
	}
	if role != `"test-e2e-role"` {
		t.Fatalf("role = %q, want the verbatim quoted segment", role)
	}
	if !wgo || aa || all {
		t.Fatalf("flags = (%v, %v, %v), want (true, false, false)", wgo, aa, all)
	}
	if !reflect.DeepEqual(privs, []string{"USAGE"}) {
		t.Fatalf("privileges = %q, want [USAGE]", privs)
	}
}

// TestBackfillGrantPrivilegesFields_StrictPrivilegeManagementPinScope locks
// the scope of the strict_privilege_management pin. The vendored v2.19.0
// schema has the attribute on grant_privileges_to_account_role only. The
// backfill must pin it there, and must not touch it on database-role
// grants. d.Set on a missing key makes the SDK log "[ERROR] setting state:
// Invalid address to set" on every Read. That log spam is what the
// roleNameKey gate removes. The rest of the backfill must still run for
// database roles.
func TestBackfillGrantPrivilegesFields_StrictPrivilegeManagementPinScope(t *testing.T) {
	newID := func() string { return "parent-minimal-role|true|false|CREATE DATABASE|OnAccount" }
	captureErrorLogs := func(t *testing.T) *bytes.Buffer {
		t.Helper()
		var buf bytes.Buffer
		log.SetOutput(&buf)
		t.Cleanup(func() { log.SetOutput(os.Stderr) })
		return &buf
	}

	t.Run("account_role pins the field its schema has", func(t *testing.T) {
		logs := captureErrorLogs(t)
		sc := testGrantPrivilegesSchema(accountRoleNameKey, true)
		d := schema.TestResourceDataRaw(t, sc, map[string]any{})
		d.SetId(newID())

		backfillGrantPrivilegesFields(d, accountRoleNameKey)

		if got := d.Get("strict_privilege_management").(bool); got {
			t.Fatalf("strict_privilege_management = true, want pinned false")
		}
		if s := logs.String(); strings.Contains(s, "Invalid address to set") {
			t.Fatalf("unexpected SDK state errors on the account-role path: %q", s)
		}
	})

	t.Run("database_role never touches the field its schema lacks", func(t *testing.T) {
		logs := captureErrorLogs(t)
		sc := testGrantPrivilegesSchema(databaseRoleNameKey, false)
		d := schema.TestResourceDataRaw(t, sc, map[string]any{})
		d.SetId(newID())

		// Control: prove the capture actually sees SDK state-set errors on
		// this schema. Without it, a broken capture could make the absence
		// assertion below pass vacuously.
		_ = d.Set("strict_privilege_management", false)
		if s := logs.String(); !strings.Contains(s, "Invalid address to set") {
			t.Fatalf("log capture control failed - capture is broken: %q", s)
		}
		logs.Reset()

		backfillGrantPrivilegesFields(d, databaseRoleNameKey)

		if s := logs.String(); strings.Contains(s, "Invalid address to set") {
			t.Fatalf("SDK state error logged for missing key: %q", s)
		}
		// The gate must not skip the rest of the backfill.
		if got := d.Get("with_grant_option").(bool); !got {
			t.Fatalf("with_grant_option = false, want true backfilled from the ID")
		}
		if got := d.Get(databaseRoleNameKey).(string); got != "parent-minimal-role" {
			t.Fatalf("database_role_name = %q, want %q", got, "parent-minimal-role")
		}
	})
}
