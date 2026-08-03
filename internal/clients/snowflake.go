package clients

import (
	"context"
	"encoding/json"

	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
	tfsdk "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/upjet/v2/pkg/terraform"

	clusterv1beta1 "github.com/Cube-Asia/provider-upjet-snowflake/apis/cluster/v1beta1"
	namespacedv1beta1 "github.com/Cube-Asia/provider-upjet-snowflake/apis/namespaced/v1beta1"
)

const (
	// error messages
	errNoProviderConfig     = "no providerConfigRef provided"
	errGetProviderConfig    = "cannot get referenced ProviderConfig"
	errTrackUsage           = "cannot track ProviderConfig usage"
	errExtractCredentials   = "cannot extract credentials"
	errUnmarshalCredentials = "cannot unmarshal snowflake credentials as JSON"
	errConfigureProvider    = "cannot configure the Snowflake terraform provider"
)

// For the full list of supported config keys, see the Snowflake TF provider schema:
// https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#schema

// Config keys matching the Snowflake TF provider schema.
// See https://registry.terraform.io/providers/snowflakedb/snowflake/latest/docs#schema
const (
	keyAccountName                   = "account_name"
	keyOrganizationName              = "organization_name"
	keyUser                          = "user"
	keyPassword                      = "password"
	keyAuthenticator                 = "authenticator"
	keyPrivateKey                    = "private_key"
	keyPrivateKeyPassphrase          = "private_key_passphrase"
	keyToken                         = "token"
	keyRole                          = "role"
	keyWarehouse                     = "warehouse"
	keyHost                          = "host"
	keyProtocol                      = "protocol"
	keyPort                          = "port"
	keyProfile                       = "profile"
	keyPasscode                      = "passcode"
	keyPasscodeInPassword            = "passcode_in_password"
	keyOktaURL                       = "okta_url"
	keyLoginTimeout                  = "login_timeout"
	keyRequestTimeout                = "request_timeout"
	keyClientTimeout                 = "client_timeout"
	keyJwtClientTimeout              = "jwt_client_timeout"
	keyJwtExpireTimeout              = "jwt_expire_timeout"
	keyExternalBrowserTimeout        = "external_browser_timeout"
	keyMaxRetryCount                 = "max_retry_count"
	keyClientRequestMfaToken         = "client_request_mfa_token"
	keyClientStoreTemporaryCred      = "client_store_temporary_credential"
	keyKeepSessionAlive              = "keep_session_alive"
	keyValidateDefaultParameters     = "validate_default_parameters"
	keyOcspFailOpen                  = "ocsp_fail_open"
	keyDisableOcspChecks             = "disable_ocsp_checks"
	keyDisableQueryContextCache      = "disable_query_context_cache"
	keyInsecureMode                  = "insecure_mode"
	keyDisableTelemetry              = "disable_telemetry"
	keyIncludeRetryReason            = "include_retry_reason"
	keyDisableConsoleLogin           = "disable_console_login"
	keyDisableSamlURLCheck           = "disable_saml_url_check"
	keyTmpDirPath                    = "tmp_dir_path"
	keyDriverTracing                 = "driver_tracing"
	keyOauthClientID                 = "oauth_client_id"
	keyOauthClientSecret             = "oauth_client_secret"
	keyOauthAuthorizationURL         = "oauth_authorization_url"
	keyOauthTokenRequestURL          = "oauth_token_request_url"
	keyOauthRedirectURI              = "oauth_redirect_uri"
	keyOauthScope                    = "oauth_scope"
	keyWorkloadIdentityProvider      = "workload_identity_provider"
	keyWorkloadIdentityEntraResource = "workload_identity_entra_resource"
	keyEnableSingleUseRefreshTokens  = "enable_single_use_refresh_tokens"
	keyCertRevocationCheckMode       = "cert_revocation_check_mode"
	keyCrlAllowCertsWithoutCrlURL    = "crl_allow_certificates_without_crl_url"
	keyCrlInMemoryCacheDisabled      = "crl_in_memory_cache_disabled"
	keyCrlOnDiskCacheDisabled        = "crl_on_disk_cache_disabled"
	keyCrlHTTPClientTimeout          = "crl_http_client_timeout"
	keyProxyHost                     = "proxy_host"
	keyProxyPort                     = "proxy_port"
	keyProxyUser                     = "proxy_user"
	keyProxyPassword                 = "proxy_password"
	keyProxyProtocol                 = "proxy_protocol"
	keyNoProxy                       = "no_proxy"
	keyLogQueryText                  = "log_query_text"
	keyLogQueryParameters            = "log_query_parameters"
	keySkipTomlFilePermVerification  = "skip_toml_file_permission_verification"
	keyUseLegacyTomlFile             = "use_legacy_toml_file"
	keyParams                        = "params"
)

// configKeys lists all supported Snowflake provider config keys.
var configKeys = []string{
	keyAccountName,
	keyOrganizationName,
	keyUser,
	keyPassword,
	keyAuthenticator,
	keyPrivateKey,
	keyPrivateKeyPassphrase,
	keyToken,
	keyRole,
	keyWarehouse,
	keyHost,
	keyProtocol,
	keyPort,
	keyProfile,
	keyPasscode,
	keyPasscodeInPassword,
	keyOktaURL,
	keyLoginTimeout,
	keyRequestTimeout,
	keyClientTimeout,
	keyJwtClientTimeout,
	keyJwtExpireTimeout,
	keyExternalBrowserTimeout,
	keyMaxRetryCount,
	keyClientRequestMfaToken,
	keyClientStoreTemporaryCred,
	keyKeepSessionAlive,
	keyValidateDefaultParameters,
	keyOcspFailOpen,
	keyDisableOcspChecks,
	keyDisableQueryContextCache,
	keyInsecureMode,
	keyDisableTelemetry,
	keyIncludeRetryReason,
	keyDisableConsoleLogin,
	keyDisableSamlURLCheck,
	keyTmpDirPath,
	keyDriverTracing,
	keyOauthClientID,
	keyOauthClientSecret,
	keyOauthAuthorizationURL,
	keyOauthTokenRequestURL,
	keyOauthRedirectURI,
	keyOauthScope,
	keyWorkloadIdentityProvider,
	keyWorkloadIdentityEntraResource,
	keyEnableSingleUseRefreshTokens,
	keyCertRevocationCheckMode,
	keyCrlAllowCertsWithoutCrlURL,
	keyCrlInMemoryCacheDisabled,
	keyCrlOnDiskCacheDisabled,
	keyCrlHTTPClientTimeout,
	keyProxyHost,
	keyProxyPort,
	keyProxyUser,
	keyProxyPassword,
	keyProxyProtocol,
	keyNoProxy,
	keyLogQueryText,
	keyLogQueryParameters,
	keySkipTomlFilePermVerification,
	keyUseLegacyTomlFile,
	keyParams,
}

// buildProviderConfiguration translates the extracted credential map into the
// terraform.Setup configuration map, JSON-decoding the params key.
func buildProviderConfiguration(creds map[string]string) (map[string]any, error) {
	cfg := make(map[string]any, len(configKeys))
	for _, k := range configKeys {
		v, ok := creds[k]
		if !ok {
			continue
		}
		if k != keyParams {
			cfg[k] = v
			continue
		}
		var params map[string]string
		if err := json.Unmarshal([]byte(v), &params); err != nil {
			return nil, errors.Wrap(err, "cannot unmarshal params as JSON map")
		}
		cfg[k] = params
	}
	return cfg, nil
}

// TerraformSetupBuilder builds Terraform a terraform.SetupFn function which
// returns Terraform provider setup configuration
func TerraformSetupBuilder(version, providerSource, providerVersion string, ujprovider *ujconfig.Provider) terraform.SetupFn {
	return func(ctx context.Context, client client.Client, mg resource.Managed) (terraform.Setup, error) {
		ps := terraform.Setup{
			Version: version,
			Requirement: terraform.ProviderRequirement{
				Source:  providerSource,
				Version: providerVersion,
			},
		}

		pcSpec, err := resolveProviderConfig(ctx, client, mg)
		if err != nil {
			return terraform.Setup{}, errors.Wrap(err, "cannot resolve provider config")
		}

		data, err := resource.CommonCredentialExtractor(ctx, pcSpec.Credentials.Source, client, pcSpec.Credentials.CommonCredentialSelectors)
		if err != nil {
			return ps, errors.Wrap(err, errExtractCredentials)
		}
		creds := map[string]string{}
		if err := json.Unmarshal(data, &creds); err != nil {
			return ps, errors.Wrap(err, errUnmarshalCredentials)
		}

		ps.Configuration, err = buildProviderConfiguration(creds)
		if err != nil {
			return ps, err
		}

		if ujprovider == nil || ujprovider.TerraformProvider == nil {
			return ps, errors.New(errConfigureProvider + ": no terraform provider configured")
		}
		diags := ujprovider.TerraformProvider.Configure(ctx, &tfsdk.ResourceConfig{Config: ps.Configuration})
		if diags.HasError() {
			return ps, errors.Errorf("%s: %v", errConfigureProvider, diags)
		}
		ps.Meta = ujprovider.TerraformProvider.Meta()
		return ps, nil
	}
}

func toSharedPCSpec(pc *clusterv1beta1.ProviderConfig) (*namespacedv1beta1.ProviderConfigSpec, error) {
	if pc == nil {
		return nil, nil
	}
	data, err := json.Marshal(pc.Spec)
	if err != nil {
		return nil, err
	}

	var mSpec namespacedv1beta1.ProviderConfigSpec
	err = json.Unmarshal(data, &mSpec)
	return &mSpec, err
}

func resolveProviderConfig(ctx context.Context, crClient client.Client, mg resource.Managed) (*namespacedv1beta1.ProviderConfigSpec, error) {
	switch managed := mg.(type) {
	case resource.LegacyManaged: //nolint:staticcheck // still handling cluster-scoped behavior
		return resolveLegacy(ctx, crClient, managed)
	case resource.ModernManaged:
		return resolveModern(ctx, crClient, managed)
	default:
		return nil, errors.New("resource is not a managed resource")
	}
}

func resolveLegacy(ctx context.Context, client client.Client, mg resource.LegacyManaged) (*namespacedv1beta1.ProviderConfigSpec, error) { //nolint:staticcheck // still handling cluster-scoped behavior
	configRef := mg.GetProviderConfigReference()
	if configRef == nil {
		return nil, errors.New(errNoProviderConfig)
	}
	pc := &clusterv1beta1.ProviderConfig{}
	if err := client.Get(ctx, types.NamespacedName{Name: configRef.Name}, pc); err != nil {
		return nil, errors.Wrap(err, errGetProviderConfig)
	}

	t := resource.NewLegacyProviderConfigUsageTracker(client, &clusterv1beta1.ProviderConfigUsage{})
	if err := t.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, errTrackUsage)
	}

	return toSharedPCSpec(pc)
}

func resolveModern(ctx context.Context, crClient client.Client, mg resource.ModernManaged) (*namespacedv1beta1.ProviderConfigSpec, error) {
	configRef := mg.GetProviderConfigReference()
	if configRef == nil {
		return nil, errors.New(errNoProviderConfig)
	}

	pcRuntimeObj, err := crClient.Scheme().New(namespacedv1beta1.SchemeGroupVersion.WithKind(configRef.Kind))
	if err != nil {
		return nil, errors.Wrap(err, "unknown GVK for ProviderConfig")
	}
	pcObj, ok := pcRuntimeObj.(client.Object)
	if !ok {
		// This indicates a programming error, types are not properly generated
		return nil, errors.New(" is not an Object")
	}

	// Namespace will be ignored if the PC is a cluster-scoped type
	if err := crClient.Get(ctx, types.NamespacedName{Name: configRef.Name, Namespace: mg.GetNamespace()}, pcObj); err != nil {
		return nil, errors.Wrap(err, errGetProviderConfig)
	}

	var pcSpec namespacedv1beta1.ProviderConfigSpec
	pcu := &namespacedv1beta1.ProviderConfigUsage{}
	switch pc := pcObj.(type) {
	case *namespacedv1beta1.ProviderConfig:
		pcSpec = pc.Spec
		if pcSpec.Credentials.SecretRef != nil {
			pcSpec.Credentials.SecretRef.Namespace = mg.GetNamespace()
		}
	case *namespacedv1beta1.ClusterProviderConfig:
		pcSpec = pc.Spec
	default:
		return nil, errors.New("unknown provider config type")
	}
	t := resource.NewProviderConfigUsageTracker(crClient, pcu)
	if err := t.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, errTrackUsage)
	}
	return &pcSpec, nil
}
