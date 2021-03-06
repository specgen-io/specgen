package gents

import (
	"github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/specgen-io/specgen/v2/genopenapi"
)

func GenerateExpressService(serviceFile string, swaggerPath string, generatePath string, servicesPath string, validation string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}

	sources := []gen.TextFile{}

	for _, version := range specification.Versions {
		sources = append(sources, *generateServiceApis(&version, generatePath))
		sources = append(sources, *generateExpressVersionRouting(&version, validation, generatePath))
	}
	sources = append(sources, *generateExpressSpecRouter(specification, generatePath))

	modelsFiles := generateModels(specification, validation, generatePath)
	sources = append(sources, modelsFiles...)

	sources = append(sources, *genopenapi.GenerateOpenapi(specification, swaggerPath))

	if servicesPath != "" {
		services := generateServicesImplementations(specification, servicesPath)
		err = gen.WriteFiles(services, false)
		if err != nil {
			return err
		}
	}

	err = gen.WriteFiles(sources, true)
	if err != nil {
		return err
	}

	return nil
}

func versionFilename(version *spec.Version, filename string, ext string) string {
	if version.Version.Source != "" {
		filename = filename + "_" + version.Version.FlatCase()
	}
	if ext != "" {
		filename = filename + "." + ext
	}
	return filename
}

func versionModule(version *spec.Version, filename string) string {
	if version.Version.Source != "" {
		filename = filename + "_" + version.Version.FlatCase()
	}
	return filename
}
