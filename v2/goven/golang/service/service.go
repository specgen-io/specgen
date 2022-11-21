package service

import (
	"github.com/specgen-io/specgen/v2/goven/generator"
	"github.com/specgen-io/specgen/v2/goven/openapi"
	"github.com/specgen-io/specgen/v2/goven/spec"
)

func GenerateService(specification *spec.Spec, moduleName string, swaggerPath string, generatePath string, servicesPath string) *generator.Sources {
	sources := generator.NewSources()

	modules := NewModules(moduleName, generatePath, servicesPath, specification)
	generator := NewGenerator(modules)

	sources.AddGenerated(generator.RootRouting(specification))
	sources.AddGeneratedAll(generator.AllStaticFiles())
	sources.AddGenerated(generator.ErrorModels(specification.HttpErrors))
	sources.AddGeneratedAll(generator.HttpErrors(&specification.HttpErrors.Responses))
	for _, version := range specification.Versions {
		sources.AddGeneratedAll(generator.Routings(&version))
		sources.AddGeneratedAll(generator.ServicesInterfaces(&version))
		sources.AddGenerated(generator.Models(&version))
	}

	if swaggerPath != "" {
		sources.AddGenerated(openapi.GenerateOpenapi(specification, swaggerPath))
	}

	if servicesPath != "" {
		for _, version := range specification.Versions {
			sources.AddScaffoldedAll(generator.ServicesImpls(&version))
		}
	}

	return sources
}
