package config

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// ---------------------------------------------------------------------------
// External name pattern helpers
//
// Snowflake's Terraform provider uses two ID encoding functions:
//
//   EncodeResourceIdentifier  — newer, used by stable resources.
//     AccountObjectIdentifier → Name() (bare, no quotes)
//     DatabaseObjectIdentifier → FullyQualifiedName() → "db"."name"
//     SchemaObjectIdentifier → FullyQualifiedName() → "db"."schema"."name"
//     Multiple parts join with |.
//
//   EncodeSnowflakeID  — legacy, still used by preview resources.
//     AccountObjectIdentifier → Name() (bare, no quotes)
//     DatabaseObjectIdentifier → db|name (unquoted, pipe-separated parts)
//     SchemaObjectIdentifier → db|schema|name (unquoted, pipe-separated parts)
//
// The helpers below encode the corresponding template for each pattern.
// See .work/snowflakedb/snowflake/pkg/helpers/ for the canonical implementation.
// ---------------------------------------------------------------------------

// SchemaObjectIdentifier returns external name config for Snowflake's
// 3-part SchemaObjectIdentifier fully-qualified-name format:
//
//	"db"."schema"."name"
//
// The name is omitted from spec.forProvider (set via metadata.name →
// crossplane.io/external-name annotation).
//
// Used by stable-family resources that use EncodeResourceIdentifier
// with a SchemaObjectIdentifier argument — the ID is produced by
// FullyQualifiedName() which quotes each segment.
func SchemaObjectIdentifier() config.ExternalName {
	return config.TemplatedStringAsIdentifier("name",
		`"{{ .parameters.database }}"."{{ .parameters.schema }}"."{{ .external_name }}"`)
}

// PipeSeparatedIdentifier returns external name config for Snowflake's
// pipe-separated 3-part ID format:
//
//	db|schema|name
//
// Used by preview-family resources that use EncodeSnowflakeID
// with a SchemaObjectIdentifier argument — the ID is produced by
// joining databaseName, schemaName, and Name() with |.
func PipeSeparatedIdentifier() config.ExternalName {
	return config.TemplatedStringAsIdentifier("name",
		`{{ .parameters.database }}|{{ .parameters.schema }}|{{ .external_name }}`)
}

// DatabaseSchemaIdentifier returns external name config for Snowflake's
// 2-part quoted ID format:
//
//	"db"."name"
//
// Used by resources whose ID is a DatabaseObjectIdentifier
// (e.g. schema, database_role).
func DatabaseSchemaIdentifier() config.ExternalName {
	return config.TemplatedStringAsIdentifier("name",
		`"{{ .parameters.database }}"."{{ .external_name }}"`)
}

// ---------------------------------------------------------------------------
// External name configurations
// ---------------------------------------------------------------------------

var ExternalNameConfigs = map[string]config.ExternalName{
	// =========================================================================
	// Stable resources (ShortGroup: stable)
	// =========================================================================

	// snowflake_account: import with '"<org>"."<account>"' — compound, but
	// the external name is just the bare account name. Org comes from the
	// provider config. AccountObjectIdentifier via EncodeResourceIdentifier.
	"snowflake_account": config.NameAsIdentifier,

	// snowflake_account_parameter: import with '<key>' — ID is the parameter
	// key. No 'name' field; uses 'key' via ParameterAsIdentifier.
	"snowflake_account_parameter": config.ParameterAsIdentifier("key"),

	// snowflake_account_role: import with '"<account_role_name>"' — bare name.
	"snowflake_account_role": config.NameAsIdentifier,

	// snowflake_account_session_policy_attachment: import format is the
	// fully qualified session policy name (FQN). The FQN is stored in the
	// session_policy_name field which doubles as the identifier.
	"snowflake_account_session_policy_attachment": config.ParameterAsIdentifier("session_policy_name"),

	// snowflake_api_authentication_integration_*: integration names.
	"snowflake_api_authentication_integration_with_authorization_code_grant": config.NameAsIdentifier,
	"snowflake_api_authentication_integration_with_client_credentials":       config.NameAsIdentifier,
	"snowflake_api_authentication_integration_with_jwt_bearer":               config.NameAsIdentifier,

	// snowflake_authentication_policy: SchemaObjectIdentifier.
	"snowflake_authentication_policy": SchemaObjectIdentifier(),

	// snowflake_catalog_integration_*: integration names.
	"snowflake_catalog_integration_aws_glue":       config.NameAsIdentifier,
	"snowflake_catalog_integration_iceberg_rest":   config.NameAsIdentifier,
	"snowflake_catalog_integration_object_storage": config.NameAsIdentifier,
	"snowflake_catalog_integration_open_catalog":   config.NameAsIdentifier,

	// snowflake_compute_pool: bare name.
	"snowflake_compute_pool": config.NameAsIdentifier,

	// snowflake_current_account: singleton — the literal string "current_account".
	"snowflake_current_account": config.IdentifierFromProvider,

	// snowflake_current_organization_account: import with '"<name>"' — bare name.
	"snowflake_current_organization_account": config.NameAsIdentifier,

	// snowflake_database: bare name (AccountObjectIdentifier).
	"snowflake_database": config.NameAsIdentifier,

	// snowflake_database_role: DatabaseObjectIdentifier — "db"."role".
	"snowflake_database_role": DatabaseSchemaIdentifier(),

	// snowflake_execute: random UUID generated by the TF provider.
	"snowflake_execute": config.IdentifierFromProvider,

	// snowflake_external_oauth_integration: bare name.
	"snowflake_external_oauth_integration": config.NameAsIdentifier,

	// snowflake_external_volume: bare name.
	"snowflake_external_volume": config.NameAsIdentifier,

	// snowflake_git_repository: SchemaObjectIdentifier.
	"snowflake_git_repository": SchemaObjectIdentifier(),

	// snowflake_grant_account_role: compound with conditional grantee type.
	// ID: '<role_name>|ROLE|<parent_role_name>' or
	//     '<role_name>|USER|<user_name>'.
	"snowflake_grant_account_role": config.TemplatedStringAsIdentifier(
		"role_name",
		"\"{{ .external_name }}\"|{{ if .parameters.parent_role_name }}ROLE|\"{{ .parameters.parent_role_name }}\"{{ else }}USER|\"{{ .parameters.user_name }}\"{{ end }}",
	),

	// snowflake_grant_application_role: compound with conditional grantee type.
	// ID: '<app_role_fqn>|ACCOUNT_ROLE|<parent_account_role>' or
	//     '<app_role_fqn>|APPLICATION|<application>'.
	"snowflake_grant_application_role": config.TemplatedStringAsIdentifier(
		"application_role_name",
		"{{ .external_name }}|{{ if .parameters.parent_account_role_name }}ACCOUNT_ROLE|\"{{ .parameters.parent_account_role_name }}\"{{ else }}APPLICATION|\"{{ .parameters.application_name }}\"{{ end }}",
	),

	// snowflake_grant_database_role: compound with conditional grantee type.
	// ID: '<db_role_fqn>|ROLE|<parent_role>' or
	//     '<db_role_fqn>|DATABASE_ROLE|<parent_db_role>' or
	//     '<db_role_fqn>|SHARE|<share>'.
	"snowflake_grant_database_role": config.TemplatedStringAsIdentifier(
		"database_role_name",
		"{{ .external_name }}|{{ if .parameters.parent_role_name }}ROLE|\"{{ .parameters.parent_role_name }}\"{{ else if .parameters.parent_database_role_name }}DATABASE_ROLE|\"{{ .parameters.parent_database_role_name }}\"{{ else }}SHARE|\"{{ .parameters.share_name }}\"{{ end }}",
	),

	// snowflake_grant_ownership: compound ID with 5-7+ variable parts
	// (<target_role>|<fqn>|<outbound>|<kind>|<object_data>).
	// ponytail: IdentifierFromProvider — revisit when upjet supports
	// dynamic-length template IDs.
	"snowflake_grant_ownership": config.IdentifierFromProvider,

	// snowflake_grant_privileges_to_account_role: compound ID with 5-9
	// variable parts (<role>|<with_grant>|<always>|<privileges>|<kind>|<data>).
	// ponytail: IdentifierFromProvider — revisit when upjet supports
	// dynamic-length template IDs.
	"snowflake_grant_privileges_to_account_role": config.IdentifierFromProvider,

	// snowflake_grant_privileges_to_database_role: compound ID with 6-9
	// variable parts (<db_role>|<with_grant>|<always>|<privileges>|<kind>|<data>).
	// ponytail: IdentifierFromProvider — revisit when upjet supports
	// dynamic-length template IDs.
	"snowflake_grant_privileges_to_database_role": config.IdentifierFromProvider,

	// snowflake_grant_privileges_to_share: compound ID with conditional target.
	// ID: '<share>|<privileges>|OnDatabase|<db>' or OnSchema/OnTable/etc.
	"snowflake_grant_privileges_to_share": config.TemplatedStringAsIdentifier(
		"to_share",
		"\"{{ .external_name }}\"|{{ .parameters.privileges }}|{{ if .parameters.on_database }}OnDatabase|\"{{ .parameters.on_database }}\"{{ else if .parameters.on_schema }}OnSchema|\"{{ .parameters.on_schema }}\"{{ else if .parameters.on_table }}OnTable|\"{{ .parameters.on_table }}\"{{ else if .parameters.on_all_tables_in_schema }}OnAllTablesInSchema|\"{{ .parameters.on_all_tables_in_schema }}\"{{ else if .parameters.on_tag }}OnTag|\"{{ .parameters.on_tag }}\"{{ else if .parameters.on_view }}OnView|\"{{ .parameters.on_view }}\"{{ end }}",
	),

	// snowflake_image_repository: SchemaObjectIdentifier.
	"snowflake_image_repository": SchemaObjectIdentifier(),

	// snowflake_legacy_service_user: same as user — bare name.
	"snowflake_legacy_service_user": config.NameAsIdentifier,

	// snowflake_listing: bare name.
	"snowflake_listing": config.NameAsIdentifier,

	// snowflake_masking_policy: SchemaObjectIdentifier.
	"snowflake_masking_policy": SchemaObjectIdentifier(),

	// snowflake_network_policy: bare name.
	"snowflake_network_policy": config.NameAsIdentifier,

	// snowflake_network_rule: SchemaObjectIdentifier.
	"snowflake_network_rule": SchemaObjectIdentifier(),

	// snowflake_oauth_integration_*: integration names.
	"snowflake_oauth_integration_for_custom_clients":       config.NameAsIdentifier,
	"snowflake_oauth_integration_for_partner_applications": config.NameAsIdentifier,

	// snowflake_password_policy: SchemaObjectIdentifier.
	"snowflake_password_policy": SchemaObjectIdentifier(),

	// snowflake_primary_connection: bare name.
	"snowflake_primary_connection": config.NameAsIdentifier,

	// snowflake_resource_monitor: bare name.
	"snowflake_resource_monitor": config.NameAsIdentifier,

	// snowflake_row_access_policy: SchemaObjectIdentifier.
	"snowflake_row_access_policy": SchemaObjectIdentifier(),

	// snowflake_saml2_integration: bare name.
	"snowflake_saml2_integration": config.NameAsIdentifier,

	// snowflake_schema: DatabaseObjectIdentifier — "db"."schema".
	"snowflake_schema": DatabaseSchemaIdentifier(),

	// snowflake_scim_integration: bare name.
	"snowflake_scim_integration": config.NameAsIdentifier,

	// snowflake_secondary_connection: bare name.
	"snowflake_secondary_connection": config.NameAsIdentifier,

	// snowflake_secondary_database: bare name.
	"snowflake_secondary_database": config.NameAsIdentifier,

	// snowflake_secret_with_*: SchemaObjectIdentifier (all variants).
	"snowflake_secret_with_authorization_code_grant": SchemaObjectIdentifier(),
	"snowflake_secret_with_basic_authentication":     SchemaObjectIdentifier(),
	"snowflake_secret_with_client_credentials":       SchemaObjectIdentifier(),
	"snowflake_secret_with_generic_string":           SchemaObjectIdentifier(),

	// snowflake_service: SchemaObjectIdentifier.
	"snowflake_service": SchemaObjectIdentifier(),

	// snowflake_service_user: same as user — bare name.
	"snowflake_service_user": config.NameAsIdentifier,

	// snowflake_session_policy: SchemaObjectIdentifier.
	"snowflake_session_policy": SchemaObjectIdentifier(),

	// snowflake_shared_database: bare name.
	"snowflake_shared_database": config.NameAsIdentifier,

	// snowflake_stage_external_*: SchemaObjectIdentifier (all variants).
	"snowflake_stage_external_azure":         SchemaObjectIdentifier(),
	"snowflake_stage_external_gcs":           SchemaObjectIdentifier(),
	"snowflake_stage_external_s3":            SchemaObjectIdentifier(),
	"snowflake_stage_external_s3_compatible": SchemaObjectIdentifier(),

	// snowflake_stage_internal: SchemaObjectIdentifier.
	"snowflake_stage_internal": SchemaObjectIdentifier(),

	// snowflake_storage_integration_*: integration names (provider-specific).
	"snowflake_storage_integration_aws":   config.NameAsIdentifier,
	"snowflake_storage_integration_azure": config.NameAsIdentifier,
	"snowflake_storage_integration_gcs":   config.NameAsIdentifier,

	// snowflake_stream_on_*: SchemaObjectIdentifier (all variants).
	"snowflake_stream_on_directory_table": SchemaObjectIdentifier(),
	"snowflake_stream_on_external_table":  SchemaObjectIdentifier(),
	"snowflake_stream_on_table":           SchemaObjectIdentifier(),
	"snowflake_stream_on_view":            SchemaObjectIdentifier(),

	// snowflake_streamlit: SchemaObjectIdentifier.
	"snowflake_streamlit": SchemaObjectIdentifier(),

	// snowflake_tag: SchemaObjectIdentifier.
	"snowflake_tag": SchemaObjectIdentifier(),

	// snowflake_tag_association: compound ID
	// (TAG_DB.TAG_SCHEMA.TAG_NAME|TAG_VALUE|OBJECT_TYPE).
	// ponytail: IdentifierFromProvider — varies by object_type.
	// Revisit when upjet supports dynamic-length template IDs.
	"snowflake_tag_association": config.IdentifierFromProvider,

	// snowflake_task: SchemaObjectIdentifier.
	"snowflake_task": SchemaObjectIdentifier(),

	// snowflake_user: bare name.
	"snowflake_user": config.NameAsIdentifier,

	// snowflake_user_programmatic_access_token: ID is '"<user>"|"<name>"' —
	// compound of user-supplied parameters (user + token name).
	"snowflake_user_programmatic_access_token": config.TemplatedStringAsIdentifier(
		"name",
		"\"{{ .parameters.user }}\"|\"{{ .external_name }}\"",
	),

	// snowflake_user_session_policy_attachment: ID is
	// '"<user>"|"<db>"."<schema>"."<session_policy>"' —
	// compound of user (external name) and session_policy_name field.
	"snowflake_user_session_policy_attachment": config.TemplatedStringAsIdentifier(
		"user_name",
		"\"{{ .external_name }}\"|{{ .parameters.session_policy_name }}",
	),

	// snowflake_view: SchemaObjectIdentifier.
	"snowflake_view": SchemaObjectIdentifier(),

	// snowflake_warehouse: bare name.
	"snowflake_warehouse": config.NameAsIdentifier,

	// =========================================================================
	// Preview resources (ShortGroup: preview)
	// =========================================================================

	// snowflake_account_authentication_policy_attachment: FQN stored in
	// authentication_policy field doubles as the identifier.
	"snowflake_account_authentication_policy_attachment": config.ParameterAsIdentifier("authentication_policy"),

	// snowflake_account_password_policy_attachment: FQN stored in
	// password_policy field doubles as the identifier.
	"snowflake_account_password_policy_attachment": config.ParameterAsIdentifier("password_policy"),

	// snowflake_alert: pipe-separated SchemaObjectIdentifier (EncodeSnowflakeID).
	"snowflake_alert": PipeSeparatedIdentifier(),

	// snowflake_api_integration: bare name.
	"snowflake_api_integration": config.NameAsIdentifier,

	// snowflake_api_integration_*: integration names (all variants).
	"snowflake_api_integration_amazon_api_gateway":          config.NameAsIdentifier,
	"snowflake_api_integration_azure_api_management":        config.NameAsIdentifier,
	"snowflake_api_integration_external_mcp_dynamic_client": config.NameAsIdentifier,
	"snowflake_api_integration_external_mcp_oauth2":         config.NameAsIdentifier,
	"snowflake_api_integration_git_repository_github_app":   config.NameAsIdentifier,
	"snowflake_api_integration_git_repository_oauth2":       config.NameAsIdentifier,
	"snowflake_api_integration_git_repository_private_link": config.NameAsIdentifier,
	"snowflake_api_integration_git_repository_token":        config.NameAsIdentifier,
	"snowflake_api_integration_google_cloud_api_gateway":    config.NameAsIdentifier,

	// snowflake_cortex_agent: SchemaObjectIdentifier (stable-style format).
	"snowflake_cortex_agent": SchemaObjectIdentifier(),

	// snowflake_cortex_search_service: pipe-separated SchemaObjectIdentifier.
	"snowflake_cortex_search_service": PipeSeparatedIdentifier(),

	// snowflake_dynamic_table: pipe-separated SchemaObjectIdentifier.
	"snowflake_dynamic_table": PipeSeparatedIdentifier(),

	// snowflake_email_notification_integration: bare name.
	"snowflake_email_notification_integration": config.NameAsIdentifier,

	// snowflake_external_function: ID includes argument signature
	// ("db"."schema"."func"(varchar,...)) — too complex for template.
	"snowflake_external_function": config.IdentifierFromProvider,

	// snowflake_external_table: pipe-separated SchemaObjectIdentifier.
	"snowflake_external_table": PipeSeparatedIdentifier(),

	// snowflake_failover_group: bare name.
	"snowflake_failover_group": config.NameAsIdentifier,

	// snowflake_file_format: pipe-separated SchemaObjectIdentifier.
	"snowflake_file_format": PipeSeparatedIdentifier(),

	// snowflake_function_*: ID includes name and argument types
	// ("db"."schema"."func"(varchar)) — signature from args block
	// makes template-based construction impractical.
	"snowflake_function_java":       config.IdentifierFromProvider,
	"snowflake_function_javascript": config.IdentifierFromProvider,
	"snowflake_function_python":     config.IdentifierFromProvider,
	"snowflake_function_scala":      config.IdentifierFromProvider,
	"snowflake_function_sql":        config.IdentifierFromProvider,

	// snowflake_iceberg_table_*: SchemaObjectIdentifier (stable-style format).
	"snowflake_iceberg_table_from_delta_files": SchemaObjectIdentifier(),
	"snowflake_iceberg_table_from_files":       SchemaObjectIdentifier(),

	// snowflake_job_service: SchemaObjectIdentifier (stable-style format).
	"snowflake_job_service": SchemaObjectIdentifier(),

	// snowflake_managed_account: bare name.
	"snowflake_managed_account": config.NameAsIdentifier,

	// snowflake_materialized_view: pipe-separated SchemaObjectIdentifier.
	"snowflake_materialized_view": PipeSeparatedIdentifier(),

	// snowflake_network_policy_attachment: ID is the network policy name
	// with '_attachment' suffix semantics — stored in network_policy_name.
	"snowflake_network_policy_attachment": config.ParameterAsIdentifier("network_policy_name"),

	// snowflake_notebook: SchemaObjectIdentifier (stable-style format).
	"snowflake_notebook": SchemaObjectIdentifier(),

	// snowflake_notification_integration: bare name.
	"snowflake_notification_integration": config.NameAsIdentifier,

	// snowflake_object_parameter: complex pipe-separated compound:
	// <key>|<object_type>|<object_identifier>.
	"snowflake_object_parameter": config.IdentifierFromProvider,

	// snowflake_pipe: pipe-separated SchemaObjectIdentifier.
	"snowflake_pipe": PipeSeparatedIdentifier(),

	// snowflake_postgres_instance: bare name.
	"snowflake_postgres_instance": config.NameAsIdentifier,

	// snowflake_procedure_*: ID includes name and argument types
	// ("db"."schema"."func"(varchar)) — same as function_* pattern.
	"snowflake_procedure_java":       config.IdentifierFromProvider,
	"snowflake_procedure_javascript": config.IdentifierFromProvider,
	"snowflake_procedure_python":     config.IdentifierFromProvider,
	"snowflake_procedure_scala":      config.IdentifierFromProvider,
	"snowflake_procedure_sql":        config.IdentifierFromProvider,

	// snowflake_semantic_view: SchemaObjectIdentifier (stable-style format).
	"snowflake_semantic_view": SchemaObjectIdentifier(),

	// snowflake_sequence: pipe-separated SchemaObjectIdentifier.
	"snowflake_sequence": PipeSeparatedIdentifier(),

	// snowflake_share: bare name.
	"snowflake_share": config.NameAsIdentifier,

	// snowflake_stage: pipe-separated SchemaObjectIdentifier.
	"snowflake_stage": PipeSeparatedIdentifier(),

	// snowflake_storage_integration: bare name.
	"snowflake_storage_integration": config.NameAsIdentifier,

	// snowflake_storage_lifecycle_policy: SchemaObjectIdentifier.
	"snowflake_storage_lifecycle_policy": SchemaObjectIdentifier(),

	// snowflake_table: pipe-separated SchemaObjectIdentifier.
	"snowflake_table": PipeSeparatedIdentifier(),

	// snowflake_table_column_masking_policy_application: complex compound of
	// table FQN, column name, and masking policy FQN.
	"snowflake_table_column_masking_policy_application": config.IdentifierFromProvider,

	// snowflake_table_constraint: uses '❄️' (snowflake) delimiter —
	// encodes constraint type, columns, and table.
	"snowflake_table_constraint": config.IdentifierFromProvider,

	// snowflake_table_storage_lifecycle_policy_attachment: two FQNs joined
	// by pipe — table FQN and storage lifecycle policy FQN.
	"snowflake_table_storage_lifecycle_policy_attachment": config.IdentifierFromProvider,

	// snowflake_user_authentication_policy_attachment: ID is
	// '"<user>"|"<db>"."<schema>"."<authentication_policy>"' —
	// compound of user (external name) and authentication_policy_name field.
	"snowflake_user_authentication_policy_attachment": config.TemplatedStringAsIdentifier(
		"user_name",
		"\"{{ .external_name }}\"|{{ .parameters.authentication_policy_name }}",
	),

	// snowflake_user_password_policy_attachment: ID is
	// '"<user>"|"<db>"."<schema>"."<password_policy>"' —
	// compound of user (external name) and password_policy_name field.
	"snowflake_user_password_policy_attachment": config.TemplatedStringAsIdentifier(
		"user_name",
		"\"{{ .external_name }}\"|{{ .parameters.password_policy_name }}",
	),

	// snowflake_user_public_keys: bare name (user name).
	"snowflake_user_public_keys": config.NameAsIdentifier,

	// snowflake_warehouse_adaptive: bare name.
	"snowflake_warehouse_adaptive": config.NameAsIdentifier,
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
