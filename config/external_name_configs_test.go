package config

import (
	"slices"
	"testing"

	"github.com/Cube-Asia/provider-upjet-snowflake/internal/resourcelist"
)

// TestExternalNameConfigsCoverResourcelist checks that the hand-maintained
// external-name table and the hand-maintained resource lists stay in step.
// Every configured resource must appear in internal/resourcelist. Every
// listed resource must have an external-name entry.
//
// ExternalNameConfigurations skips a table key that is missing from the
// provider schema. It does not report an error. A listed resource without a
// table entry silently falls back to default external-name handling. Either
// way, the failure appears only in a cluster, as an observe or create that
// cannot reconstruct IDs. This test fails at generation time instead. A
// resource name belongs in BOTH internal/resourcelist and
// ExternalNameConfigs, or in neither.
func TestExternalNameConfigsCoverResourcelist(t *testing.T) {
	want := append(slices.Clone(resourcelist.StableResources), resourcelist.PreviewResources...)
	if len(want) == 0 {
		// If the resourcelist is empty, the comparison below passes
		// without checking anything. Fail loudly instead.
		t.Fatal("internal/resourcelist is empty; is the list being generated or reset?")
	}

	got := make([]string, 0, len(ExternalNameConfigs))
	for name := range ExternalNameConfigs {
		got = append(got, name)
	}
	slices.Sort(want)
	slices.Sort(got)

	missing := []string{}
	for _, w := range want {
		if !slices.Contains(got, w) {
			missing = append(missing, w)
		}
	}
	extra := []string{}
	for _, g := range got {
		if !slices.Contains(want, g) {
			extra = append(extra, g)
		}
	}

	if len(missing) > 0 || len(extra) > 0 {
		t.Fatalf("ExternalNameConfigs and internal/resourcelist diverged:\n  in resourcelist but not configured: %v\n  configured but not in resourcelist: %v\nAdd new resources to both lists in the same change.", missing, extra)
	}
}
