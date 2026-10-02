package config

import (
	ujconfig "github.com/crossplane/upjet/v2/pkg/config"

	"github.com/Cube-Asia/provider-upjet-snowflake/internal/resourcelist"
)

// The group Configure functions below are shared by both provider scopes:
// GetProvider and GetProviderNamespaced call them unchanged. Scope
// differences live in the ujconfig.Provider options (root group, example
// manifest namespace), not in these functions; keeping a single copy is
// what stops the two scopes from drifting apart.

const (
	stableShortGroup  = "stable"
	previewShortGroup = "preview"
)

// ConfigurePreviewGroup wires every preview resource with its short group
// and PascalCase kind.
func ConfigurePreviewGroup(p *ujconfig.Provider) {
	configureGroup(p, previewShortGroup, resourcelist.PreviewResources, nil)
}

// ConfigureStableGroup wires every stable resource with its short group and
// PascalCase kind, and applies the shared user-type workaround.
func ConfigureStableGroup(p *ujconfig.Provider) {
	configureGroup(p, stableShortGroup, resourcelist.StableResources, func(name string, r *ujconfig.Resource) {
		// All user-type resources share the same userSchema with
		// BooleanDefault ("default") and IntDefault (-1) sentinel
		// values that fail the provider's own validation.
		if name == "snowflake_legacy_service_user" || name == "snowflake_service_user" || name == "snowflake_user" {
			userConfigurationInjector(r)
		}
	})
}

// configureGroup wires one resource group: each resource gets its short
// group and PascalCase kind, and any group-specific extra configuration.
// Both scopes must stay byte-for-byte consistent; this single loop is what
// keeps them from drifting apart.
func configureGroup(p *ujconfig.Provider, shortGroup string, names []string, extra func(name string, r *ujconfig.Resource)) {
	for _, name := range names {
		p.AddResourceConfigurator(name, func(r *ujconfig.Resource) {
			r.ShortGroup = shortGroup
			r.Kind = resourcelist.ToKind(name)
			if extra != nil {
				extra(name, r)
			}
		})
	}
}

// userConfigurationInjector excludes problematic user fields from
// late-initialization via IgnoredFields. These fields have Terraform schema
// Default values that cause issues when the late-initializer copies them
// from TF state into spec.forProvider:
//
//   - Sentinel defaults: BooleanDefault("default") fails ValidateBooleanString,
//     IntDefault(-1) fails IntAtLeast(0). The TF plugin applies defaults
//     AFTER validation internally, but once late-init copies them into
//     spec.forProvider, they go into main.tf.json as user-supplied values
//     that fail the SDK's ValidateFunc — creating a perpetual loop.
//   - Case normalization: Snowflake returns "IGNORE" for unsupported_ddl_action
//     while the TF schema Default is "ignore". Late-init copies the lowercase
//     default into spec, it goes into main.tf.json as "ignore", and TF sees
//     a perpetual config-vs-state diff.
//
// IgnoredFields prevents this leak by keeping these fields out of the
// late-initialization step entirely. They only enter spec.forProvider when
// explicitly set by the user, in which case they ARE written to main.tf.json
// and managed normally by Terraform (no loop because the value is valid).
func userConfigurationInjector(r *ujconfig.Resource) {
	r.LateInitializer.IgnoredFields = append(
		r.LateInitializer.IgnoredFields,
		"disable_mfa", "disabled", "must_change_password",
		"mins_to_bypass_mfa", "mins_to_unlock",
		"unsupported_ddl_action",
	)
}
