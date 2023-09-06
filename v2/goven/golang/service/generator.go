package service

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/goven/generator"
	"github.com/specgen-io/specgen/v2/goven/golang/models"
	"github.com/specgen-io/specgen/v2/goven/golang/types"
	"github.com/specgen-io/specgen/v2/goven/spec"
)

type ServerGenerator interface {
	RootRouting(specification *spec.Spec) *generator.CodeFile
	Routings(version *spec.Version) []generator.CodeFile
	GenerateUrlParamsCtor() *generator.CodeFile
}

type Generator struct {
	ServerGenerator
	models.Generator
	Types   *types.Types
	Modules *Modules
}

func NewGenerator(jsonmode, server string, modules *Modules) *Generator {
	types := types.NewTypes()
	models := models.NewGenerator(jsonmode, &(modules.Modules))

	var serverGenerator ServerGenerator = nil
	switch server {
	case Vestigo:
		serverGenerator = NewVestigoGenerator(types, models, modules)
		break
	case HttpRouter:
		serverGenerator = NewHttpRouterGenerator(types, models, modules)
		break
	case Chi:
		serverGenerator = NewChiGenerator(types, models, modules)
		break
	default:
		panic(fmt.Sprintf(`Unsupported server: %s`, server))
	}

	return &Generator{
		serverGenerator,
		models,
		types,
		modules,
	}
}
