package models

import (
	"github.com/specgen-io/specgen/v2/goven/kotlin/packages"
	"github.com/specgen-io/specgen/v2/goven/spec"
)

type Packages struct {
	models       map[string]packages.Package
	Json         packages.Package
	JsonAdapters packages.Package
	Errors       packages.Package
	ErrorsModels packages.Package
}

func NewPackages(packageName, generatePath string, specification *spec.Spec) *Packages {
	if packageName == "" {
		packageName = specification.Name.SnakeCase()
	}
	generated := packages.New(generatePath, packageName)
	json := generated.Subpackage("json")
	adapters := json.Subpackage("adapters")
	errors := generated.Subpackage("errors")
	errorsModels := errors.Subpackage("models")

	models := map[string]packages.Package{}
	for _, version := range specification.Versions {
		models[version.Name.Source] = generated.Subpackage(version.Name.FlatCase()).Subpackage("models")
	}

	return &Packages{
		models,
		json,
		adapters,
		errors,
		errorsModels,
	}
}

func (p *Packages) Models(version *spec.Version) packages.Package {
	return p.models[version.Name.Source]
}
