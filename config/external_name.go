package config

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// ExternalNameConfigs contains all external name configurations for this
// provider.
var ExternalNameConfigs = map[string]config.ExternalName{
	// =============================================================================
	// Stable User resources
	// =============================================================================

	// snowflake_user: import with '"<user_name>"' — ID is the user name.
	// ponytail: uses name field as identifier, omitted from CRD spec (set via external-name annotation).
	"snowflake_user": config.NameAsIdentifier,

	// snowflake_service_user: same pattern as snowflake_user.
	"snowflake_service_user": config.NameAsIdentifier,

	// snowflake_legacy_service_user: same pattern as snowflake_user.
	"snowflake_legacy_service_user": config.NameAsIdentifier,

	// snowflake_user_programmatic_access_token: compound ID '"<user>"|"<name>"'.
	// ID combines user and token name — cannot derive from a single field.
	"snowflake_user_programmatic_access_token": config.IdentifierFromProvider,

	// snowflake_user_session_policy_attachment: compound ID
	// '"<user_name>"|"<database>"."<schema>"."<session_policy>"'.
	"snowflake_user_session_policy_attachment": config.IdentifierFromProvider,

	// =============================================================================
	// Preview User resources (subject to breaking changes)
	// =============================================================================

	// snowflake_user_public_keys: uses name field (user name) as identifier.
	"snowflake_user_public_keys": config.NameAsIdentifier,

	// snowflake_user_authentication_policy_attachment: no import section;
	// ID format is provider-generated compound.
	"snowflake_user_authentication_policy_attachment": config.IdentifierFromProvider,

	// snowflake_user_password_policy_attachment: pipe-separated compound ID
	// "DATABASE|SCHEMA|POLICY|USER".
	"snowflake_user_password_policy_attachment": config.IdentifierFromProvider,
}

// ExternalNameConfigurations applies all external name configs listed in the
// table ExternalNameConfigs and sets the version of those resources to v1beta1
// assuming they will be tested.
func ExternalNameConfigurations() config.ResourceOption {
	return func(r *config.Resource) {
		if e, ok := ExternalNameConfigs[r.Name]; ok {
			r.ExternalName = e
		}
	}
}

// ExternalNameConfigured returns the list of all resources whose external name
// is configured manually.
func ExternalNameConfigured() []string {
	l := make([]string, len(ExternalNameConfigs))
	i := 0
	for name := range ExternalNameConfigs {
		// $ is added to match the exact string since the format is regex.
		l[i] = name + "$"
		i++
	}
	return l
}
