package client

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/goven/generator"
	"github.com/specgen-io/specgen/v2/goven/java/models"
	"github.com/specgen-io/specgen/v2/goven/java/types"
	"github.com/specgen-io/specgen/v2/goven/spec"
)

type ClientGenerator interface {
	Clients(version *spec.Version) []generator.CodeFile
	Utils(responses *spec.ErrorResponses) []generator.CodeFile
	Exceptions(errors *spec.ErrorResponses) []generator.CodeFile
}

type Generator struct {
	ClientGenerator
	models.Generator
	Jsonlib  string
	Types    *types.Types
	Packages *Packages
}

func NewGenerator(jsonlib string, client string, packages *Packages) *Generator {
	types := models.NewTypes(jsonlib)
	models := models.NewGenerator(jsonlib, &(packages.Packages))

	var clientGenerator ClientGenerator = nil
	switch client {
	case OkHttp:
		clientGenerator = NewOkHttpGenerator(types, models, packages)
		break
	case Micronaut:
		clientGenerator = NewMicronautGenerator(types, models, packages)
		break
	default:
		panic(fmt.Sprintf(`Unsupported client: %s`, client))
	}

	return &Generator{
		clientGenerator,
		models,
		jsonlib,
		types,
		packages,
	}
}
