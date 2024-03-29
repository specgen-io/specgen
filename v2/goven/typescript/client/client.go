package client

import (
	"github.com/specgen-io/specgen/v2/goven/generator"
	"github.com/specgen-io/specgen/v2/goven/spec"
)

func GenerateClient(specification *spec.Spec, generatePath string, client string, validationName string) *generator.Sources {
	sources := generator.NewSources()
	modules := NewModules(validationName, generatePath, specification)
	generator := NewClientGenerator(client, validationName, modules)
	sources.AddGenerated(generator.SetupLibrary())
	sources.AddGenerated(generator.ParamsBuilder())
	sources.AddGenerated(generator.Errors(specification.HttpErrors))
	sources.AddGenerated(generator.ErrorModels(specification.HttpErrors))
	sources.AddGenerated(generator.ErrorResponses(specification.HttpErrors))
	for _, version := range specification.Versions {
		sources.AddGenerated(generator.Models(&version))
		for _, api := range version.Http.Apis {
			sources.AddGenerated(generator.ApiClient(&api))
		}
	}
	return sources
}
