package stable

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

const shortGroup = "stable"

// Configure configures the stable resource group for namespaced scope.
// Contains Snowflake resources with subCategory "Stable".
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("snowflake_account", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_account_parameter", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_account_role", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "AccountRole"
	})
	p.AddResourceConfigurator("snowflake_account_session_policy_attachment", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "AccountSessionPolicyAttachment"
	})
	p.AddResourceConfigurator("snowflake_api_authentication_integration_with_authorization_code_grant", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_api_authentication_integration_with_client_credentials", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_api_authentication_integration_with_jwt_bearer", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_authentication_policy", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "AuthenticationPolicy"
	})
	p.AddResourceConfigurator("snowflake_catalog_integration_aws_glue", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_catalog_integration_iceberg_rest", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_catalog_integration_object_storage", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_catalog_integration_open_catalog", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_compute_pool", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_current_account", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_current_organization_account", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_database", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_database_role", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_execute", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_external_oauth_integration", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_external_volume", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_git_repository", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_grant_account_role", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "GrantAccountRole"
	})
	p.AddResourceConfigurator("snowflake_grant_application_role", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_grant_database_role", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_grant_ownership", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_grant_privileges_to_account_role", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_grant_privileges_to_database_role", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_grant_privileges_to_share", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_image_repository", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_legacy_service_user", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "LegacyServiceUser"
	})
	p.AddResourceConfigurator("snowflake_listing", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_masking_policy", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "MaskingPolicy"
	})
	p.AddResourceConfigurator("snowflake_network_policy", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_network_rule", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_oauth_integration_for_custom_clients", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_oauth_integration_for_partner_applications", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_password_policy", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "PasswordPolicy"
	})
	p.AddResourceConfigurator("snowflake_primary_connection", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_resource_monitor", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_row_access_policy", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "RowAccessPolicy"
	})
	p.AddResourceConfigurator("snowflake_saml2_integration", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_schema", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_scim_integration", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_secondary_connection", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_secondary_database", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "SecondaryDatabase"
	})
	p.AddResourceConfigurator("snowflake_secret_with_authorization_code_grant", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_secret_with_basic_authentication", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_secret_with_client_credentials", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_secret_with_generic_string", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_service", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_service_user", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "ServiceUser"
	})
	p.AddResourceConfigurator("snowflake_session_policy", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "SessionPolicy"
	})
	p.AddResourceConfigurator("snowflake_shared_database", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "SharedDatabase"
	})
	p.AddResourceConfigurator("snowflake_stage_external_azure", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_stage_external_gcs", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_stage_external_s3", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_stage_external_s3_compatible", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_stage_internal", func(r *config.Resource) {
		r.ShortGroup = shortGroup
		r.Kind = "StageInternal"
	})
	p.AddResourceConfigurator("snowflake_storage_integration_aws", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_storage_integration_azure", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_storage_integration_gcs", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_stream_on_directory_table", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_stream_on_external_table", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_stream_on_table", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_stream_on_view", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_streamlit", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_tag", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_tag_association", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_task", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_user", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_user_programmatic_access_token", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_user_session_policy_attachment", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_view", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
	p.AddResourceConfigurator("snowflake_warehouse", func(r *config.Resource) {
		r.ShortGroup = shortGroup
	})
}
