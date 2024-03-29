package generators

import (
	"github.com/specgen-io/specgen/v2/goven/generator"
	"github.com/specgen-io/specgen/v2/goven/spec"
	"github.com/specgen-io/specgen/v2/goven/typescript/client"
	"github.com/specgen-io/specgen/v2/goven/typescript/service"
	"github.com/specgen-io/specgen/v2/goven/typescript/validations"
)

var ValidationTsValues = []string{"superstruct", "io-ts"}
var ClientTsValues = []string{"axios", "node-fetch", "browser-fetch"}
var ServerTsValues = []string{"express", "koa"}

var Models = generator.Generator{
	"models-ts",
	"TypeScript Models",
	"Generate TypeScript models source code",
	[]generator.GeneratorArg{
		{Arg: generator.ArgSpecFile, Required: true},
		{Arg: generator.ArgValidation, Required: true, Values: ValidationTsValues},
		{Arg: generator.ArgGeneratePath, Required: true},
	},
	func(specification *spec.Spec, params generator.GeneratorArgsValues) *generator.Sources {
		return validations.GenerateModels(specification, params[generator.ArgValidation], params[generator.ArgGeneratePath])
	},
}

var Client = generator.Generator{
	"client-ts",
	"TypeScript Client",
	"Generate TypeScript client source code",
	[]generator.GeneratorArg{
		{Arg: generator.ArgSpecFile, Required: true},
		{Arg: generator.ArgClient, Required: true, Values: ClientTsValues},
		{Arg: generator.ArgValidation, Required: true, Values: ValidationTsValues},
		{Arg: generator.ArgGeneratePath, Required: true},
	},
	func(specification *spec.Spec, params generator.GeneratorArgsValues) *generator.Sources {
		return client.GenerateClient(specification, params[generator.ArgGeneratePath], params[generator.ArgClient], params[generator.ArgValidation])
	},
}

var Service = generator.Generator{
	"service-ts",
	"TypeScript Service",
	"Generate TypeScript client source code",
	[]generator.GeneratorArg{
		{Arg: generator.ArgSpecFile, Required: true},
		{Arg: generator.ArgServer, Required: true, Values: ServerTsValues},
		{Arg: generator.ArgValidation, Required: true, Values: ValidationTsValues},
		{Arg: generator.ArgGeneratePath, Required: true},
		{Arg: generator.ArgServicesPath, Required: false},
		{Arg: generator.ArgSwaggerPath, Required: false},
	},
	func(specification *spec.Spec, params generator.GeneratorArgsValues) *generator.Sources {
		return service.GenerateService(specification, params[generator.ArgSwaggerPath], params[generator.ArgGeneratePath], params[generator.ArgServicesPath], params[generator.ArgServer], params[generator.ArgValidation])
	},
}

var All = []generator.Generator{
	Models,
	Client,
	Service,
}
