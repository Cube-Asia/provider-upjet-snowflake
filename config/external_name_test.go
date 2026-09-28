package config

import (
	"strings"
	"testing"
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
