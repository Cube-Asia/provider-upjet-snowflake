package config

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

// Encode-side tests for the grant ID builders. The builders reconstruct the
// compound, pipe-separated import IDs that the upstream Snowflake provider's
// parsers accept on observe and import; a wrong part count or a wrong marker
// makes every observe fail (see the OnObject regression this suite grew out
// of). Each case asserts the FULL ID string, so an upstream shape change
// fails loudly here instead of in a cluster.
//
// Expected shapes are derived from the upstream encoder, not from this
// fork's own output: helpers.EncodeResourceIdentifier joins parts with the
// pipe delimiter (pkg/helpers/resource_identifier.go:18), and the ID types
// fix the part order:
//   - GrantAccountRoleId.String: roleFQN | objType | granteeFQN
//     (pkg/resources/grant_account_role_identifier.go:16)
//   - GrantDatabaseRoleId.String: dbRoleFQN | objType | granteeFQN, where
//     the parent-database-role marker is sdk.ObjectTypeDatabaseRole =
//     "DATABASE ROLE" with a space (grant_database_role_identifier.go:16,
//     pkg/sdk/object_types.go:22)
//   - GrantPrivilegesToShareId.String: shareFQN | privileges-joined | kind
//     | targetFQN (grant_privileges_to_share_identifier.go:30)
//   - GrantOwnershipId.String: targetKind | roleName | outboundPrivs (empty
//     when nil) | blockKind | dataParts (grant_ownership_identifier.go:65)
//
// The privileges base "role|wgo|alwaysApply|privs" encoding mirrors the
// three grant-privileges resources' shared ID layout.

func TestBuildGrantAccountRoleID(t *testing.T) {
	cases := map[string]struct {
		parameters map[string]any
		want       string
		wantErr    string
	}{
		"ParentRole": {
			parameters: map[string]any{
				"role_name":        `"ANALYST"`,
				"parent_role_name": `"SYSADMIN"`,
			},
			want: `"ANALYST"|ROLE|"SYSADMIN"`,
		},
		"User": {
			parameters: map[string]any{
				"role_name": `"ANALYST"`,
				"user_name": `"etl-svc"`,
			},
			want: `"ANALYST"|USER|"etl-svc"`,
		},
		"BareNamesAreQuoted": {
			parameters: map[string]any{
				"role_name":        "analyst",
				"parent_role_name": "sysadmin",
			},
			want: `"analyst"|ROLE|"sysadmin"`,
		},
		"NoGrantee": {
			parameters: map[string]any{"role_name": `"ANALYST"`},
			wantErr:    "neither parent_role_name nor user_name is set",
		},
		"MissingRole": {
			parameters: map[string]any{"user_name": `"etl-svc"`},
			wantErr:    "role_name is required",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := buildGrantAccountRoleID(tc.parameters)
			assertIDOrError(t, got, err, tc.want, tc.wantErr)
		})
	}
}

func TestBuildGrantApplicationRoleID(t *testing.T) {
	cases := map[string]struct {
		parameters map[string]any
		want       string
		wantErr    string
	}{
		"ParentAccountRole": {
			parameters: map[string]any{
				"application_role_name":    `"APP_ROLE"`,
				"parent_account_role_name": `"SYSADMIN"`,
			},
			want: `"APP_ROLE"|ACCOUNT_ROLE|"SYSADMIN"`,
		},
		"Application": {
			parameters: map[string]any{
				"application_role_name": `"APP_ROLE"`,
				"application_name":      `"MY_APP"`,
			},
			want: `"APP_ROLE"|APPLICATION|"MY_APP"`,
		},
		"NoGrantee": {
			parameters: map[string]any{"application_role_name": `"APP_ROLE"`},
			wantErr:    "neither parent_account_role_name nor application_name is set",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := buildGrantApplicationRoleID(tc.parameters)
			assertIDOrError(t, got, err, tc.want, tc.wantErr)
		})
	}
}

func TestBuildGrantDatabaseRoleID(t *testing.T) {
	cases := map[string]struct {
		parameters map[string]any
		want       string
		wantErr    string
	}{
		"ParentRole": {
			parameters: map[string]any{
				"database_role_name": `"DB"."ANALYST"`,
				"parent_role_name":   `"SYSADMIN"`,
			},
			want: `"DB"."ANALYST"|ROLE|"SYSADMIN"`,
		},
		"ParentDatabaseRole": {
			parameters: map[string]any{
				"database_role_name":        `"DB"."ANALYST"`,
				"parent_database_role_name": `"DB"."ADMIN"`,
			},
			// sdk.ObjectTypeDatabaseRole.String() is "DATABASE ROLE" with a
			// space, not "DATABASE_ROLE".
			want: `"DB"."ANALYST"|DATABASE ROLE|"DB"."ADMIN"`,
		},
		"Share": {
			parameters: map[string]any{
				"database_role_name": `"DB"."ANALYST"`,
				"share_name":         `"SALES_SHARE"`,
			},
			want: `"DB"."ANALYST"|SHARE|"SALES_SHARE"`,
		},
		"NoGrantee": {
			parameters: map[string]any{"database_role_name": `"DB"."ANALYST"`},
			wantErr:    "none of parent_role_name, parent_database_role_name, share_name is set",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := buildGrantDatabaseRoleID(tc.parameters)
			assertIDOrError(t, got, err, tc.want, tc.wantErr)
		})
	}
}

func TestGrantPrivilegesBaseStr(t *testing.T) {
	cases := map[string]struct {
		roleName   string
		parameters map[string]any
		want       string
	}{
		"PrivilegesSortedCommaSeparated": {
			roleName: `"R"`,
			parameters: map[string]any{
				"privileges": []any{"SELECT", "INSERT", "UPDATE"},
			},
			want: `"R"|false|false|INSERT,SELECT,UPDATE`,
		},
		"AllPrivilegesCollapsesTheList": {
			roleName: `"R"`,
			parameters: map[string]any{
				"all_privileges": true,
				"privileges":     []any{"SELECT"},
			},
			want: `"R"|false|false|ALL`,
		},
		"FlagsEncodeAsTrueFalse": {
			roleName: `"R"`,
			parameters: map[string]any{
				"with_grant_option": true,
				"always_apply":      true,
				"privileges":        []any{"SELECT"},
			},
			want: `"R"|true|true|SELECT`,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := grantPrivilegesBaseStr(tc.roleName, tc.parameters); got != tc.want {
				t.Errorf("grantPrivilegesBaseStr() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestBuildGrantPrivilegesToAccountRoleID(t *testing.T) {
	cases := map[string]struct {
		parameters map[string]any
		want       string
		wantErr    string
	}{
		"OnAccount": {
			parameters: map[string]any{
				"account_role_name": `"R"`,
				"privileges":        []any{"ATTACH POLICY"},
				"on_account":        true,
			},
			want: `"R"|false|false|ATTACH POLICY|OnAccount`,
		},
		"OnAccountObject": {
			parameters: map[string]any{
				"account_role_name": `"R"`,
				"privileges":        []any{"SELECT"},
				"on_account_object": []any{
					map[string]any{"object_type": "DATABASE", "object_name": `"DB"`},
				},
			},
			want: `"R"|false|false|SELECT|OnAccountObject|DATABASE|"DB"`,
		},
		"OnSchema": {
			parameters: map[string]any{
				"account_role_name": `"R"`,
				"privileges":        []any{"USAGE"},
				"on_schema": []any{
					map[string]any{"schema_name": `"DB"."PUBLIC"`},
				},
			},
			want: `"R"|false|false|USAGE|OnSchema|OnSchema|"DB"."PUBLIC"`,
		},
		"OnSchemaObjectSingle": {
			parameters: map[string]any{
				"account_role_name": `"R"`,
				"privileges":        []any{"SELECT"},
				"on_schema_object": []any{
					map[string]any{"object_type": "TABLE", "object_name": `"DB"."PUBLIC"."TBL"`},
				},
			},
			want: `"R"|false|false|SELECT|OnSchemaObject|OnObject|TABLE|"DB"."PUBLIC"."TBL"`,
		},
		"OnSchemaObjectAll": {
			parameters: map[string]any{
				"account_role_name": `"R"`,
				"privileges":        []any{"SELECT"},
				"on_schema_object": []any{
					map[string]any{
						"all": []any{
							map[string]any{"object_type_plural": "TABLES", "in_schema": `"DB"."PUBLIC"`},
						},
					},
				},
			},
			want: `"R"|false|false|SELECT|OnSchemaObject|OnAll|TABLES|InSchema|"DB"."PUBLIC"`,
		},
		"NoTarget": {
			parameters: map[string]any{
				"account_role_name": `"R"`,
				"privileges":        []any{"SELECT"},
			},
			wantErr: "one of on_account, on_account_object, on_schema, or on_schema_object is required",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := buildGrantPrivilegesToAccountRoleID(tc.parameters)
			assertIDOrError(t, got, err, tc.want, tc.wantErr)
		})
	}
}

func TestBuildGrantPrivilegesToDatabaseRoleID(t *testing.T) {
	cases := map[string]struct {
		parameters map[string]any
		want       string
		wantErr    string
	}{
		"OnDatabase": {
			parameters: map[string]any{
				"database_role_name": `"DB"."R"`,
				"privileges":         []any{"USAGE"},
				"on_database":        `"DB"`,
			},
			want: `"DB"."R"|false|false|USAGE|OnDatabase|"DB"`,
		},
		"OnSchemaFuture": {
			parameters: map[string]any{
				"database_role_name": `"DB"."R"`,
				"privileges":         []any{"SELECT"},
				"on_schema": []any{
					map[string]any{"future_schemas_in_database": `"DB"`},
				},
			},
			want: `"DB"."R"|false|false|SELECT|OnSchema|OnFutureSchemasInDatabase|"DB"`,
		},
		"OnSchemaObjectSingleWithoutAllPrivileges": {
			// The exact shape that regressed: all_privileges unset must still
			// carry the OnObject marker.
			parameters: map[string]any{
				"database_role_name": `"DB"."R"`,
				"privileges":         []any{"SELECT"},
				"on_schema_object": []any{
					map[string]any{"object_type": "TABLE", "object_name": `"DB"."PUBLIC"."TBL"`},
				},
			},
			want: `"DB"."R"|false|false|SELECT|OnSchemaObject|OnObject|TABLE|"DB"."PUBLIC"."TBL"`,
		},
		"NoTarget": {
			parameters: map[string]any{
				"database_role_name": `"DB"."R"`,
				"privileges":         []any{"SELECT"},
			},
			wantErr: "one of on_database, on_schema, or on_schema_object is required",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := buildGrantPrivilegesToDatabaseRoleID(tc.parameters)
			assertIDOrError(t, got, err, tc.want, tc.wantErr)
		})
	}
}

func TestBuildGrantPrivilegesToShareID(t *testing.T) {
	cases := map[string]struct {
		parameters map[string]any
		want       string
		wantErr    string
	}{
		"OnTable": {
			parameters: map[string]any{
				"to_share":   `"SALES_SHARE"`,
				"privileges": []any{"SELECT"},
				"on_table":   `"DB"."PUBLIC"."TBL"`,
			},
			want: `"SALES_SHARE"|SELECT|OnTable|"DB"."PUBLIC"."TBL"`,
		},
		"OnDatabase": {
			parameters: map[string]any{
				"to_share":    `"SALES_SHARE"`,
				"privileges":  []any{"USAGE"},
				"on_database": `"DB"`,
			},
			want: `"SALES_SHARE"|USAGE|OnDatabase|"DB"`,
		},
		"PrivilegesRequired": {
			parameters: map[string]any{
				"to_share": `"SALES_SHARE"`,
				"on_table": `"DB"."PUBLIC"."TBL"`,
			},
			wantErr: "privileges is required",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := buildGrantPrivilegesToShareID(tc.parameters)
			assertIDOrError(t, got, err, tc.want, tc.wantErr)
		})
	}
}

func TestBuildGrantOwnershipID(t *testing.T) {
	cases := map[string]struct {
		parameters map[string]any
		want       string
		wantErr    string
	}{
		"OnObjectToAccountRole": {
			parameters: map[string]any{
				"account_role_name":   `"R"`,
				"outbound_privileges": "COPY",
				"on": []any{
					map[string]any{"object_type": "TABLE", "object_name": `"DB"."PUBLIC"."TBL"`},
				},
			},
			want: `ToAccountRole|"R"|COPY|OnObject|TABLE|"DB"."PUBLIC"."TBL"`,
		},
		"OnAllToDatabaseRole": {
			parameters: map[string]any{
				"database_role_name": `"DB"."R"`,
				"on": []any{
					map[string]any{
						"all": []any{
							map[string]any{"object_type_plural": "TABLES", "in_database": `"DB"`},
						},
					},
				},
			},
			want: `ToDatabaseRole|"DB"."R"||OnAll|TABLES|InDatabase|"DB"`,
		},
		"NoOnBlock": {
			parameters: map[string]any{"account_role_name": `"R"`},
			wantErr:    "'on' block is required",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := buildGrantOwnershipID(tc.parameters)
			assertIDOrError(t, got, err, tc.want, tc.wantErr)
		})
	}
}

func TestNormalizeSFObjectID(t *testing.T) {
	cases := map[string]struct{ in, want string }{
		"Bare":   {"analyst", `"analyst"`},
		"Quoted": {`"analyst"`, `"analyst"`},
		"Fully":  {`"db"."analyst"`, `"db"."analyst"`},
		"Mixed":  {`db."analyst"`, `"db"."analyst"`},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := normalizeSFObjectID(tc.in); got != tc.want {
				t.Errorf("normalizeSFObjectID(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

// assertIDOrError keeps the tables small: each case declares either the
// wanted ID or a wanted error substring.
func assertIDOrError(t *testing.T, got string, err error, wantID, wantErr string) {
	t.Helper()
	if wantErr != "" {
		if err == nil {
			t.Fatalf("want error %q, got ID %q", wantErr, got)
		}
		if !strings.Contains(err.Error(), wantErr) {
			t.Fatalf("error %q does not contain %q", err.Error(), wantErr)
		}
		return
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != wantID {
		t.Errorf("ID = %q, want %q", got, wantID)
	}
}

// ---------------------------------------------------------------------------
// Round-trip and determinism (adopted from the standalone round-trip PR)
// ---------------------------------------------------------------------------

// The grant-privileges external ID is the storage format for everything the
// Read path cannot recover from Snowflake: the role name, grant option, and
// privilege set all survive only inside this string. The upstream builders
// emit it bare, while IDs observed in a live Snowflake account are quoted by
// Snowflake's ID delegation ("role"|true|...), so the parser must accept
// quoted segments and the encoder/parser pair must round-trip losslessly.
// If either side changes shape, grants owned by the old format stop being
// recognized.

func TestGrantPrivilegesAccountRoleIDRoundTrip(t *testing.T) {
	params := map[string]any{
		"account_role_name": "test-e2e-role",
		"with_grant_option": true,
		"privileges":        []any{"USAGE"},
		"on_account_object": []any{map[string]any{"object_type": "DATABASE", "object_name": "test-e2e-db"}},
	}

	id, err := buildGrantPrivilegesToAccountRoleID(params)
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	want := "test-e2e-role|true|false|USAGE|OnAccountObject|DATABASE|test-e2e-db"
	if id != want {
		t.Fatalf("builder id = %q, want %q", id, want)
	}

	role, wgo, aa, all, privs, ok := parseGrantPrivilegesBaseID(id)
	if !ok {
		t.Fatalf("parse(%q) failed", id)
	}
	if role != "test-e2e-role" || !wgo || aa || all {
		t.Fatalf("parse(%q) = (%q, %v, %v, %v)", id, role, wgo, aa, all)
	}
	if !reflect.DeepEqual(privs, []string{"USAGE"}) {
		t.Fatalf("privileges = %q, want [USAGE]", privs)
	}

	// Rebuilding from the parsed fields must reproduce the exact ID.
	// (The builder consumes the JSON-decoded []any shape; the parser
	// returns []string, so convert.)
	privsAny := make([]any, len(privs))
	for i, p := range privs {
		privsAny[i] = p
	}
	rebuilt := map[string]any{
		"account_role_name": role,
		"with_grant_option": wgo,
		"always_apply":      aa,
		"all_privileges":    all,
		"privileges":        privsAny,
		"on_account_object": []any{map[string]any{"object_type": "DATABASE", "object_name": "test-e2e-db"}},
	}
	id2, err := buildGrantPrivilegesToAccountRoleID(rebuilt)
	if err != nil {
		t.Fatalf("rebuild: %v", err)
	}
	if id2 != id {
		t.Fatalf("rebuild id = %q, want original %q", id2, id)
	}
}

func TestGrantPrivilegesAccountRoleIDDeterministicPrivilegeOrder(t *testing.T) {
	build := func(privs []any) (string, error) {
		return buildGrantPrivilegesToAccountRoleID(map[string]any{
			"account_role_name": "r",
			"with_grant_option": false,
			"privileges":        privs,
			"on_account_object": []any{map[string]any{"object_type": "DATABASE", "object_name": "d"}},
		})
	}
	// The same set in a different input order must yield the same ID —
	// otherwise Create and Read would disagree about which grant exists.
	idAB, err := build([]any{"CREATE USER", "CREATE DATABASE"})
	if err != nil {
		t.Fatalf("build AB: %v", err)
	}
	idBA, err := build([]any{"CREATE DATABASE", "CREATE USER"})
	if err != nil {
		t.Fatalf("build BA: %v", err)
	}
	if idAB != idBA {
		t.Fatalf("privilege order changed the ID: %q vs %q", idAB, idBA)
	}
	_, _, _, _, privs, ok := parseGrantPrivilegesBaseID(idAB)
	if !ok {
		t.Fatalf("parse(%q) failed", idAB)
	}
	got := append([]string(nil), privs...)
	sort.Strings(got)
	if !reflect.DeepEqual(got, []string{"CREATE DATABASE", "CREATE USER"}) {
		t.Fatalf("privileges = %q, want the input set", privs)
	}
}

func TestGrantPrivilegesDatabaseRoleIDRoundTrip(t *testing.T) {
	params := map[string]any{
		"database_role_name": "db-rw",
		"with_grant_option":  false,
		"privileges":         []any{"USAGE"},
		"on_database":        "analytics",
	}
	id, err := buildGrantPrivilegesToDatabaseRoleID(params)
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	// Pin the exact wire form for the bare-name shape — the table above
	// covers the quoted form ("DB"), so together both segment styles are
	// drift-pinned without duplicating an assertion.
	want := "db-rw|false|false|USAGE|OnDatabase|analytics"
	if id != want {
		t.Fatalf("builder id = %q, want %q", id, want)
	}
	role, wgo, aa, all, privs, ok := parseGrantPrivilegesBaseID(id)
	if !ok {
		t.Fatalf("parse(%q) failed", id)
	}
	if role != "db-rw" || wgo || aa || all {
		t.Fatalf("parse(%q) = (%q, %v, %v, %v)", id, role, wgo, aa, all)
	}
	if !reflect.DeepEqual(privs, []string{"USAGE"}) {
		t.Fatalf("privileges = %q, want [USAGE]", privs)
	}
}
