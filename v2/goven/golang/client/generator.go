package client

import (
	"github.com/specgen-io/specgen/v2/goven/generator"
	"github.com/specgen-io/specgen/v2/goven/golang/empty"
	"github.com/specgen-io/specgen/v2/goven/golang/models"
	"github.com/specgen-io/specgen/v2/goven/golang/types"
	"github.com/specgen-io/specgen/v2/goven/spec"
)

type ClientGenerator interface {
	Clients(version *spec.Version) []generator.CodeFile
	ErrorsHandler(errors spec.ErrorResponses) *generator.CodeFile
	ResponseHelperFunctions() *generator.CodeFile
}

type Generator struct {
	models.Generator
	ClientGenerator
	Types   *types.Types
	Modules *Modules
}

func NewGenerator(jsonmode string, modules *Modules) *Generator {
	types := types.NewTypes()
	return &Generator{
		models.NewGenerator(jsonmode, &(modules.Modules)),
		NewNetHttpGenerator(modules, types),
		types,
		modules,
	}
}

func (g *Generator) EmptyType() *generator.CodeFile {
	return empty.GenerateEmpty(g.Modules.Empty)
}

func (g *Generator) AllStaticFiles() []generator.CodeFile {
	return []generator.CodeFile{
		*g.EnumsHelperFunctions(),
		*g.EmptyType(),
		*g.TypeConverter(),
		*g.Params(),
		*g.FormDataParams(),
		*g.ResponseHelperFunctions(),
	}
}
