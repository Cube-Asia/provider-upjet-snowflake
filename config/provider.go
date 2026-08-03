package config

import (
	// Note(turkenh): we are importing this to embed provider schema document
	_ "embed"

	sfprovider "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider"
	ujconfig "github.com/crossplane/upjet/v2/pkg/config"

	previewCluster "github.com/Cube-Asia/provider-upjet-snowflake/config/cluster/preview"
	stableCluster "github.com/Cube-Asia/provider-upjet-snowflake/config/cluster/stable"
	previewNamespaced "github.com/Cube-Asia/provider-upjet-snowflake/config/namespaced/preview"
	stableNamespaced "github.com/Cube-Asia/provider-upjet-snowflake/config/namespaced/stable"
)

const (
	resourcePrefix = "snowflake"
	modulePath     = "github.com/Cube-Asia/provider-upjet-snowflake"
)

//go:embed schema.json
var providerSchema string

//go:embed provider-metadata.yaml
var providerMetadata string

// GetProvider returns provider configuration
func GetProvider() *ujconfig.Provider {
	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithRootGroup("snowflake.crossplane.io"),
		// all resources are reconciled via the Terraform Plugin SDK (no-fork
		// mode); nothing goes through the CLI-fork include list.
		ujconfig.WithIncludeList([]string{}),
		ujconfig.WithTerraformPluginSDKIncludeList(ExternalNameConfigured()),
		ujconfig.WithTerraformProvider(sfprovider.Provider()),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithDefaultResourceOptions(
			ExternalNameConfigurations(),
			GrantPrivilegesReadGapWorkaround(),
		))

	for _, configure := range []func(provider *ujconfig.Provider){
		// add custom config functions
		stableCluster.Configure,
		previewCluster.Configure,
	} {
		configure(pc)
	}

	pc.ConfigureResources()
	return pc
}

// GetProviderNamespaced returns the namespaced provider configuration
func GetProviderNamespaced() *ujconfig.Provider {
	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithRootGroup("snowflake.m.crossplane.io"),
		ujconfig.WithIncludeList([]string{}),
		ujconfig.WithTerraformPluginSDKIncludeList(ExternalNameConfigured()),
		ujconfig.WithTerraformProvider(sfprovider.Provider()),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithDefaultResourceOptions(
			ExternalNameConfigurations(),
			GrantPrivilegesReadGapWorkaround(),
		),
		ujconfig.WithExampleManifestConfiguration(ujconfig.ExampleManifestConfiguration{
			ManagedResourceNamespace: "crossplane-system",
		}))

	for _, configure := range []func(provider *ujconfig.Provider){
		// add custom config functions
		stableNamespaced.Configure,
		previewNamespaced.Configure,
	} {
		configure(pc)
	}

	pc.ConfigureResources()
	return pc
}
