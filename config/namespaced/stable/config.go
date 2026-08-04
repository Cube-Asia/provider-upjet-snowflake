package stable

import (
	ujconfig "github.com/crossplane/upjet/v2/pkg/config"

	"github.com/Cube-Asia/provider-upjet-snowflake/internal/resourcelist"
)

const shortGroup = "stable"

// Configure configures the stable resource group for namespaced scope.
func Configure(p *ujconfig.Provider) {
	for _, name := range resourcelist.StableResources {
		p.AddResourceConfigurator(name, func(r *ujconfig.Resource) {
			r.ShortGroup = shortGroup
			r.Kind = resourcelist.ToKind(name)
			// All user-type resources share the same userSchema with
			// BooleanDefault ("default") and IntDefault (-1) sentinel
			// values that fail the provider's own validation.
			if name == "snowflake_legacy_service_user" || name == "snowflake_service_user" || name == "snowflake_user" {
				userConfigurationInjector(r)
			}
		})
	}
}

func userConfigurationInjector(r *ujconfig.Resource) {
	// Exclude problematic fields from late-initialization via IgnoredFields.
	// These fields have Terraform schema Default values that cause issues when
	// the late-initializer copies them from TF state into spec.forProvider:
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
	r.LateInitializer.IgnoredFields = append(
		r.LateInitializer.IgnoredFields,
		"disable_mfa", "disabled", "must_change_password",
		"mins_to_bypass_mfa", "mins_to_unlock",
		"unsupported_ddl_action",
	)
}
