package userpreview

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

const shortGroup = "userpreview"

// Configure configures the preview user resource group (cluster-scoped).
// Contains resources with subcategory "Preview" that are subject to breaking
// changes in upstream Terraform provider releases.
//
// ponytail: isolated ShortGroup prevents API group churn from affecting stable
// user resources when preview resources change. Upstream breaking changes only
// touch the userpreview.* CRDs, never user.*.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("snowflake_user_public_keys", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_user_authentication_policy_attachment", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_user_password_policy_attachment", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
}
