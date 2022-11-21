package validations

import (
	"github.com/specgen-io/specgen/v2/goven/generator"
	"github.com/specgen-io/specgen/v2/goven/spec"
	"github.com/specgen-io/specgen/v2/goven/typescript/validations/modules"
)

func GenerateModels(specification *spec.Spec, validation string, generatePath string) *generator.Sources {
	sources := generator.NewSources()
	modules := modules.NewModules(validation, generatePath, specification)
	generator := New(validation, modules)
	sources.AddGenerated(generator.SetupLibrary())
	for _, version := range specification.Versions {
		sources.AddGenerated(generator.Models(&version))
	}
	return sources
}
