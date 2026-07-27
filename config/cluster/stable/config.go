package stable

import (
	ujconfig "github.com/crossplane/upjet/v2/pkg/config"

	"github.com/Cube-Asia/provider-upjet-snowflake/internal/resourcelist"
)

const shortGroup = "stable"

// Configure configures the stable resource group for cluster scope.
func Configure(p *ujconfig.Provider) {
	for _, name := range resourcelist.StableResources {
		p.AddResourceConfigurator(name, func(r *ujconfig.Resource) {
			r.ShortGroup = shortGroup
			r.Kind = resourcelist.ToKind(name)
		})
	}
}
