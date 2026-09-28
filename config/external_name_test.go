package config

import (
	"context"
	"strings"
	"testing"

	sfresources "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/resources"
)

// Regression coverage for the OnObject marker gate: the single-object case of
// on_schema_object must always emit the 8-part
// "…|OnSchemaObject|OnObject|<object_type>|<object_name>" ID. Omitting the
// marker when all_privileges=false produced a 7-part ID that
// ParseGrantPrivilegesTo{Account,Database}RoleId rejects with
// "invalid OnSchemaObjectGrantKind" on every observe, so the grant could
// never become Ready (upstream's OnSchemaObjectGrantData.String always
// writes OnObject, regardless of all_privileges).
func TestProcessOnSchemaObjectID_SingleObjectAlwaysEmitsOnObject(t *testing.T) {
	const base = `"R"|false|false|SELECT`

	cases := map[string]struct {
		parameters map[string]any
		want       string
	}{
		"SingleObjectWithoutAllPrivileges": {
			parameters: map[string]any{
				"on_schema_object": []any{
					map[string]any{
						"object_type": "TABLE",
						"object_name": `"DB"."PUBLIC"."TBL"`,
					},
				},
			},
			want: base + `|OnSchemaObject|OnObject|TABLE|"DB"."PUBLIC"."TBL"`,
		},
		"SingleObjectWithAllPrivileges": {
			parameters: map[string]any{
				"all_privileges": true,
				"on_schema_object": []any{
					map[string]any{
						"object_type": "TABLE",
						"object_name": `"DB"."PUBLIC"."TBL"`,
					},
				},
			},
			want: base + `|OnSchemaObject|OnObject|TABLE|"DB"."PUBLIC"."TBL"`,
		},
		"OnAllBlockUnchanged": {
			parameters: map[string]any{
				"on_schema_object": []any{
					map[string]any{
						"all": []any{
							map[string]any{
								"object_type_plural": "TABLES",
								"in_schema":          `"DB"."PUBLIC"`,
							},
						},
					},
				},
			},
			want: base + `|OnSchemaObject|OnAll|TABLES|InSchema|"DB"."PUBLIC"`,
		},
		"OnFutureBlockUnchanged": {
			parameters: map[string]any{
				"on_schema_object": []any{
					map[string]any{
						"future": []any{
							map[string]any{
								"object_type_plural": "VIEWS",
								"in_database":        `"DB"`,
							},
						},
					},
				},
			},
			want: base + `|OnSchemaObject|OnFuture|VIEWS|InDatabase|"DB"`,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, found, err := processOnSchemaObjectID(tc.parameters, base, "grant_privileges_to_database_role")
			if err != nil {
				t.Fatalf("processOnSchemaObjectID() error = %v", err)
			}
			if !found {
				t.Fatal("processOnSchemaObjectID() found = false, want true")
			}
			if got != tc.want {
				t.Errorf("processOnSchemaObjectID() = %q, want %q", got, tc.want)
			}
			if strings.Count(got, "|") != strings.Count(tc.want, "|") {
				t.Errorf("part count mismatch: got %q, want %q", got, tc.want)
			}
		})
	}
}

// TestGrantPrivilegesToDatabaseRole_GetIDFnHealsLegacyAnnotation locks the
// observe-path heal for the live stuck grants: MRs created by provider builds before
// the OnObject fix carry a stored external-name in the legacy 7-part form
// (OnSchemaObject|<object_type>|<object_name>, no marker), which
// sfresources.ParseGrantPrivilegesToDatabaseRoleId rejects with "invalid
// OnSchemaObjectGrantKind" on every observe - the grant could never go
// Ready and the spec's shadow rehearsal and smoke test cannot run.
//
// The heal is structural, not a data migration: upjet's connect step calls
// GetIDFn BEFORE anything parses an ID (upjet
// pkg/controller/external_tfpluginsdk.go), and this resource's GetIDFn
// rebuilds the ID from the complete spec parameters, falling back to the
// stored annotation only when the builder errors. So once the fixed builder
// deploys, the legacy annotation is never parsed: the rebuilt modern ID is
// what upstream's parser receives, and the stuck MR heals on its next
// reconcile. This test pins that chain end to end - including that the
// legacy string really is unparsable (the bug) and that the rebuilt string
// parses AND re-encodes identically (the heal).
func TestGrantPrivilegesToDatabaseRole_GetIDFnHealsLegacyAnnotation(t *testing.T) {
	const legacy = `"APP_DB"."APP_DB_RO_REPORTING"|false|false|SELECT|OnSchemaObject|TABLE|"APP_DB"."PUBLIC"."APP_TABLE"`
	const modern = `"APP_DB"."APP_DB_RO_REPORTING"|false|false|SELECT|OnSchemaObject|OnObject|TABLE|"APP_DB"."PUBLIC"."APP_TABLE"`

	// The claim parameters that exposed the bug, expressed as TF parameters.
	parameters := map[string]any{
		"database_role_name": `"APP_DB"."APP_DB_RO_REPORTING"`,
		"with_grant_option":  false,
		"always_apply":       false,
		"privileges":         []any{"SELECT"},
		"on_schema_object": []any{
			map[string]any{
				"object_type": "TABLE",
				"object_name": `"APP_DB"."PUBLIC"."APP_TABLE"`,
			},
		},
	}

	// The stored legacy annotation is exactly what upstream's parser rejects.
	if _, err := sfresources.ParseGrantPrivilegesToDatabaseRoleId(legacy); err == nil {
		t.Fatal("legacy ID parsed cleanly, want 'invalid OnSchemaObjectGrantKind' error")
	}

	got, err := GrantPrivilegesToDatabaseRoleIdentifier().GetIDFn(context.Background(), legacy, parameters, nil)
	if err != nil {
		t.Fatalf("GetIDFn() error = %v", err)
	}
	if got != modern {
		t.Fatalf("GetIDFn() = %q, want %q", got, modern)
	}

	// The rebuilt ID must satisfy the parser the observe path actually runs.
	id, err := sfresources.ParseGrantPrivilegesToDatabaseRoleId(got)
	if err != nil {
		t.Fatalf("ParseGrantPrivilegesToDatabaseRoleId(%q) error = %v", got, err)
	}
	// Re-encode so a silently-lossy parse cannot pass.
	if reencoded := id.String(); reencoded != modern {
		t.Fatalf("parsed ID re-encodes as %q, want %q", reencoded, modern)
	}
}
