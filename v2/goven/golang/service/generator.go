package service

import (
	"github.com/specgen-io/specgen/v2/goven/generator"
	"github.com/specgen-io/specgen/v2/goven/golang/empty"
	"github.com/specgen-io/specgen/v2/goven/golang/models"
	"github.com/specgen-io/specgen/v2/goven/golang/types"
	"github.com/specgen-io/specgen/v2/goven/spec"
)

type ServiceGenerator interface {
	RootRouting(specification *spec.Spec) *generator.CodeFile
	HttpErrors(responses *spec.Responses) []generator.CodeFile
	CheckContentType() *generator.CodeFile
	Routings(version *spec.Version) []generator.CodeFile
	ResponseHelperFunctions() *generator.CodeFile
}

type Generator struct {
	models.Generator
	ServiceGenerator
	Types   *types.Types
	Modules *Modules
}

func NewGenerator(modules *Modules) *Generator {
	modelsGenerator := models.NewGenerator(&(modules.Modules))
	types := types.NewTypes()
	return &Generator{
		modelsGenerator,
		NewVestigoGenerator(types, modelsGenerator, modules),
		types,
		modules,
	}
}

func (g *Generator) AllStaticFiles() []generator.CodeFile {
	return []generator.CodeFile{
		*g.EnumsHelperFunctions(),
		*empty.GenerateEmpty(g.Modules.Empty),
		*generateParamsParser(g.Modules.ParamsParser),
		*g.ResponseHelperFunctions(),
		*g.CheckContentType(),
	}
}
