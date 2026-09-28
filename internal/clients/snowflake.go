package clients

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"sync"
	"time"

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
	// One cache per SetupFn (there are two setups in main.go, one per scope,
	// each backed by its own *schema.Provider). It retains the configured
	// provider meta so a reconcile re-mints a Snowflake SDK session only when
	// the provider configuration changes or the entry expires, not on every
	// call.
	var cache = metaCache{ttl: time.Hour}

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

		key, err := hashConfiguration(ps.Configuration)
		if err != nil {
			return ps, errors.Wrap(err, "cannot hash provider configuration")
		}

		meta, err := cache.metaFor(key, func() (any, error) {
			return configureAndMeta(ctx, ujprovider, ps.Configuration)
		})
		if err != nil {
			return ps, err
		}
		ps.Meta = meta
		return ps, nil
	}
}

// metaCache memoizes the configured Terraform provider meta, keyed by the
// provider configuration digest.
//
// Why it exists: ConfigureProvider on the shared schema.Provider mints a
// live Snowflake session on every call. sdk.NewClient runs sqlx.Connect
// (sql.Open plus Ping) and then CurrentAccount and CurrentSession round
// trips. Nothing ever closes the previous session. Upjet invokes the
// SetupFn once per managed-resource reconcile, so before this cache every
// reconcile opened an authenticated session, a database/sql pool (with its
// connectionOpener goroutine) and an HTTP transport, then dropped the
// reference. In practice this caused a straight-line memory climb of about
// 1.5-2 GiB per hour until the pod ran out of memory and was killed.
//
// Trade-off: when the configuration (for example the private key) changes,
// the old session's meta is replaced, not closed. The SDK Client type lives
// in an internal package we cannot name here, and closing it immediately
// could remove the connection pool under concurrent reconciles that still
// hold the previous meta. A rotation therefore leaks exactly one session,
// once per change, instead of one per reconcile.
type metaCache struct {
	mu      sync.Mutex
	entries map[string]cacheEntry
	ttl     time.Duration
	// now is injectable for tests; nil means time.Now.
	now func() time.Time
}

// cacheEntry pairs a configured meta with its mint time. Entries expire
// after the cache TTL. gosnowflake renews an expired *session* token
// transparently (390112 → renewExpiredSessionToken, restful.go), but the
// *master* token has no verified recovery path in the vendored driver
// (renewRestfulSession renews with the master token itself, and a pooled
// conn idled past its lifetime fails renewal with a plain error that
// database/sql does not discard). Expiring entries bounds re-logins to
// 24 per day per configuration and removes the whole stale-session class.
type cacheEntry struct {
	meta    any
	created time.Time
}

// metaFor returns the cached meta for key, and calls configure exactly once
// on a miss. The lock is held across the check, the configure, and the
// store, so concurrent misses on a cold cache cannot race into duplicate
// sessions. Entries are keyed by the provider configuration digest, so
// reconciles that reference different ProviderConfigs (for example rotated
// credentials) keep independent sessions instead of thrashing one slot. The
// map is bounded by the number of distinct configurations this provider
// instance resolves, and each entry is at most one TTL old.
//
// The vendored provider.Context is opaque (internal package), so nothing
// here can or should Close it. Expiry only orphans the old client, which is
// reclaimed when the provider restarts.
func (c *metaCache) metaFor(key string, configure func() (any, error)) (any, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.now
	if now == nil {
		now = time.Now
	}
	if entry, ok := c.entries[key]; ok && now().Sub(entry.created) < c.ttl {
		return entry.meta, nil
	}
	meta, err := configure()
	if err != nil {
		// Do not cache failures: a transient configure error must not pin
		// the cache to a missing entry, and the next call retries.
		return nil, err
	}
	if c.entries == nil {
		c.entries = make(map[string]cacheEntry)
	}
	c.entries[key] = cacheEntry{meta: meta, created: now()}
	return meta, nil
}

// hashConfiguration derives the cache key from the provider configuration.
// json.Marshal sorts map keys, so the digest is stable across reconciles.
// The digest covers secret material (password, private key) by design, so
// it must never be logged.
func hashConfiguration(cfg map[string]any) (string, error) {
	b, err := json.Marshal(cfg)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return string(sum[:]), nil
}

// configureAndMeta runs the Terraform provider's Configure and returns the
// resulting meta. It is a package variable so tests can observe how many
// times the (session-minting) configure step runs without contacting a
// Snowflake endpoint.
var configureAndMeta = func(ctx context.Context, ujprovider *ujconfig.Provider, cfg map[string]any) (any, error) {
	if ujprovider == nil || ujprovider.TerraformProvider == nil {
		return nil, errors.New(errConfigureProvider + ": no terraform provider configured")
	}
	diags := ujprovider.TerraformProvider.Configure(ctx, &tfsdk.ResourceConfig{Config: cfg})
	if diags.HasError() {
		return nil, errors.Errorf("%s: %v", errConfigureProvider, diags)
	}
	return ujprovider.TerraformProvider.Meta(), nil
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
