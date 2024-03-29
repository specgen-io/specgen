package service

import (
	"github.com/specgen-io/specgen/v2/goven/kotlin/models"
	"github.com/specgen-io/specgen/v2/goven/kotlin/packages"
	"github.com/specgen-io/specgen/v2/goven/spec"
)

type Packages struct {
	models.Packages
	versions        map[string]VersionPackages
	Root            packages.Package
	ContentType     packages.Package
	Converters      packages.Package
	RootControllers packages.Package
}

type VersionPackages struct {
	Services     packages.Package
	Controllers  packages.Package
	ServicesImpl packages.Package
}

func NewPackages(packageName, generatePath, servicesPath string, specification *spec.Spec) *Packages {
	if packageName == "" {
		packageName = specification.Name.SnakeCase()
	}

	generated := packages.New(generatePath, packageName)
	contenttype := generated.Subpackage("contenttype")
	converters := generated.Subpackage("converters")
	controllers := generated.Subpackage("controllers")
	implementations := packages.New(servicesPath, packageName)

	versions := map[string]VersionPackages{}
	for _, version := range specification.Versions {
		main := generated.Subpackage(version.Name.FlatCase())
		services := main.Subpackage("services")
		controllers := main.Subpackage("controllers")
		servicesImpl := implementations.Subpackage("services").Subpackage(version.Name.FlatCase())

		versionPackages := VersionPackages{
			services,
			controllers,
			servicesImpl,
		}

		versions[version.Name.Source] = versionPackages
	}

	return &Packages{
		*models.NewPackages(packageName, generatePath, specification),
		versions,
		generated,
		contenttype,
		converters,
		controllers,
	}
}

func (p *Packages) ServicesApi(api *spec.Api) packages.Package {
	return p.versions[api.InHttp.InVersion.Name.Source].Services.Subpackage(api.Name.SnakeCase())
}

func (p *Packages) ServicesImpl(version *spec.Version) packages.Package {
	return p.versions[version.Name.Source].ServicesImpl
}

func (p *Packages) Controllers(version *spec.Version) packages.Package {
	return p.versions[version.Name.Source].Controllers
}
