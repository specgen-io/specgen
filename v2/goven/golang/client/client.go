package client

import (
	"github.com/specgen-io/specgen/v2/goven/generator"
	"github.com/specgen-io/specgen/v2/goven/spec"
)

func GenerateClient(specification *spec.Spec, jsonmode string, moduleName string, generatePath string) *generator.Sources {
	sources := generator.NewSources()

	modules := NewModules(moduleName, generatePath, specification)
	generator := NewGenerator(jsonmode, modules)

	sources.AddGeneratedAll(generator.AllStaticFiles())

	sources.AddGeneratedAll(generator.ErrorModels(specification.HttpErrors))
	sources.AddGenerated(generator.Errors(&specification.HttpErrors.Responses))
	sources.AddGenerated(generator.ErrorsHandler(specification.HttpErrors.Responses))

	for _, version := range specification.Versions {
		sources.AddGeneratedAll(generator.Models(&version))
		sources.AddGeneratedAll(generator.Clients(&version))
	}
	return sources
}
