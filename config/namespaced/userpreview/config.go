package userpreview

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

const shortGroup = "userpreview"

// Configure configures the preview user resource group (namespaced-scoped).
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
