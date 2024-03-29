package client

import (
	"github.com/specgen-io/specgen/v2/goven/spec"
	"github.com/specgen-io/specgen/v2/goven/typescript/module"
	"github.com/specgen-io/specgen/v2/goven/typescript/validations/modules"
)

type Modules struct {
	modules.Modules
	clients  map[string]map[string]module.Module
	Params   module.Module
	Response module.Module
}

func NewModules(validationName, generatePath string, specification *spec.Spec) *Modules {
	root := module.New(generatePath)
	params := root.Submodule("params")
	response := root.Submodule("response")
	clients := map[string]map[string]module.Module{}
	for _, version := range specification.Versions {
		clients[version.Name.Source] = map[string]module.Module{}
		versionModule := root.Submodule(version.Name.FlatCase())
		for _, api := range version.Http.Apis {
			clients[version.Name.Source][api.Name.Source] = versionModule.Submodule(api.Name.SnakeCase())

		}
	}

	return &Modules{
		*modules.NewModules(validationName, generatePath, specification),
		clients,
		params,
		response,
	}
}

func (m *Modules) Client(api *spec.Api) module.Module {
	return m.clients[api.InHttp.InVersion.Name.Source][api.Name.Source]
}
