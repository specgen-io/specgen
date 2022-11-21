package models

import (
	"github.com/specgen-io/specgen/v2/goven/generator"
	"github.com/specgen-io/specgen/v2/goven/golang/types"
	"github.com/specgen-io/specgen/v2/goven/spec"
)

type Generator interface {
	Models(version *spec.Version) *generator.CodeFile
	ErrorModels(httperrors *spec.HttpErrors) *generator.CodeFile
	EnumValuesStrings(model *spec.NamedModel) string
	EnumsHelperFunctions() *generator.CodeFile
}

func NewGenerator(modules *Modules) Generator {
	types := types.NewTypes()
	return NewEncodingJsonGenerator(types, modules)
}
