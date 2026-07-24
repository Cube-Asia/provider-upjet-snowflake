package user

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

const shortGroup = "user"

// Configure configures the stable user resource group (namespaced-scoped).
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("snowflake_user", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_service_user", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "ServiceUser"
	})
	p.AddResourceConfigurator("snowflake_legacy_service_user", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "LegacyServiceUser"
	})
	p.AddResourceConfigurator("snowflake_user_programmatic_access_token", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_user_session_policy_attachment", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
}
