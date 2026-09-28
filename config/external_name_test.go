package config

import (
	"context"
	"strings"
	"testing"

	sfresources "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/resources"
)

// onSchemaObjectKey is the shared parameters map key for the on_schema_object
// grant parameter. Both tests in this file use it. goconst (min-occurrences: 5
// on this layer) requires the repeated key to be a named constant. The
// test-file exclusion arrives with a later layer.
const onSchemaObjectKey = "on_schema_object"

// Regression coverage for the OnObject marker gate. The single-object case of
// on_schema_object must always emit the 8-part
// "…|OnSchemaObject|OnObject|<object_type>|<object_name>" ID. Omitting the
// marker when all_privileges=false produced a 7-part ID.
// ParseGrantPrivilegesTo{Account,Database}RoleId rejects that ID with
// "invalid OnSchemaObjectGrantKind" on every observe, so the grant could never
// become Ready. Upstream's OnSchemaObjectGrantData.String always writes
// OnObject, regardless of all_privileges.
func TestProcessOnSchemaObjectID_SingleObjectAlwaysEmitsOnObject(t *testing.T) {
	const base = `"R"|false|false|SELECT`

	cases := map[string]struct {
		parameters map[string]any
		want       string
	}{
		"SingleObjectWithoutAllPrivileges": {
			parameters: map[string]any{
				onSchemaObjectKey: []any{
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
				onSchemaObjectKey: []any{
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
				onSchemaObjectKey: []any{
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
				onSchemaObjectKey: []any{
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
// observe-path heal for the live stuck grants. MRs created by provider builds
// before the OnObject fix carry a stored external-name in the legacy 7-part
// form (OnSchemaObject|<object_type>|<object_name>, no marker).
// sfresources.ParseGrantPrivilegesToDatabaseRoleId rejects that form with
// "invalid OnSchemaObjectGrantKind" on every observe. The grant could never go
// Ready, and the spec's shadow rehearsal and smoke test cannot run.
//
// The heal is structural, not a data migration. Upjet's connect step calls
// GetIDFn before anything parses an ID (upjet
// pkg/controller/external_tfpluginsdk.go). This resource's GetIDFn rebuilds
// the ID from the complete spec parameters. It falls back to the stored
// annotation only when the builder errors. Once the fixed builder deploys, the
// legacy annotation is never parsed. The rebuilt modern ID is what upstream's
// parser receives, and the stuck MR heals on its next reconcile. This test
// pins that chain end to end. It proves the legacy string is unparsable (the
// bug) and that the rebuilt string parses and re-encodes identically (the
// heal).
func TestGrantPrivilegesToDatabaseRole_GetIDFnHealsLegacyAnnotation(t *testing.T) {
	const legacy = `"APP_DB"."APP_DB_RO_REPORTING"|false|false|SELECT|OnSchemaObject|TABLE|"APP_DB"."PUBLIC"."APP_TABLE"`
	const modern = `"APP_DB"."APP_DB_RO_REPORTING"|false|false|SELECT|OnSchemaObject|OnObject|TABLE|"APP_DB"."PUBLIC"."APP_TABLE"`

	// The claim parameters that exposed the bug, expressed as TF parameters.
	parameters := map[string]any{
		"database_role_name": `"APP_DB"."APP_DB_RO_REPORTING"`,
		"with_grant_option":  false,
		"always_apply":       false,
		"privileges":         []any{"SELECT"},
		onSchemaObjectKey: []any{
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
	// Re-encode so a parse that silently drops parts cannot pass.
	if reencoded := id.String(); reencoded != modern {
		t.Fatalf("parsed ID re-encodes as %q, want %q", reencoded, modern)
	}
}
