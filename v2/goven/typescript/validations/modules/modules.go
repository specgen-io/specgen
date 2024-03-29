package modules

import (
	"github.com/specgen-io/specgen/v2/goven/spec"
	"github.com/specgen-io/specgen/v2/goven/typescript/module"
)

type Modules struct {
	models       map[string]module.Module
	Errors       module.Module
	ErrorsModels module.Module
	Validation   module.Module
}

func NewModules(validationName, generatePath string, specification *spec.Spec) *Modules {
	root := module.New(generatePath)
	errors := root.SubmoduleIndex("errors")
	errorModels := errors.Submodule("models")
	validation := root.Submodule(validationName)

	models := map[string]module.Module{}
	for _, version := range specification.Versions {
		versionModule := root.Submodule(version.Name.FlatCase())
		models[version.Name.Source] = versionModule.Submodule("models")
	}

	return &Modules{
		models,
		errors,
		errorModels,
		validation,
	}
}

func (m *Modules) Models(version *spec.Version) module.Module {
	return m.models[version.Name.Source]
}
