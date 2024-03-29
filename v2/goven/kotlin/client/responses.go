package client

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/goven/generator"
	"github.com/specgen-io/specgen/v2/goven/kotlin/packages"
	"github.com/specgen-io/specgen/v2/goven/kotlin/types"
	"github.com/specgen-io/specgen/v2/goven/kotlin/writer"
	"github.com/specgen-io/specgen/v2/goven/spec"
)

func responseCreate(response *spec.OperationResponse, resultVar string) string {
	if len(response.Operation.Responses.Success()) > 1 {
		return fmt.Sprintf(`return %s.%s(%s)`, responseInterfaceName(response.Operation), response.Name.PascalCase(), resultVar)
	} else {
		if resultVar != "" {
			return fmt.Sprintf(`return %s`, resultVar)
		} else {
			return `return`
		}
	}
}

func responses(api *spec.Api, types *types.Types, apiPackage packages.Package, modelsVersionPackage packages.Package, errorModelsPackage packages.Package) []generator.CodeFile {
	files := []generator.CodeFile{}
	for _, operation := range api.Operations {
		if len(operation.Responses.Success()) > 1 {
			files = append(files, *responseInterface(types, &operation, apiPackage, modelsVersionPackage, errorModelsPackage))
		}
	}
	return files
}

func responseInterface(types *types.Types, operation *spec.NamedOperation, apiPackage packages.Package, modelsVersionPackage packages.Package, errorModelsPackage packages.Package) *generator.CodeFile {
	w := writer.New(apiPackage, responseInterfaceName(operation))
	w.Line(`import %s`, modelsVersionPackage.PackageStar)
	w.Line(`import %s`, errorModelsPackage.PackageStar)
	w.EmptyLine()
	w.Line(`interface [[.ClassName]] {`)
	for index, response := range operation.Responses.Success() {
		if index > 0 {
			w.EmptyLine()
		}
		responseImpl(w.Indented(), types, response)
	}
	w.Line(`}`)
	return w.ToCodeFile()
}

func responseImpl(w *writer.Writer, types *types.Types, response *spec.OperationResponse) {
	serviceResponseImplementationName := response.Name.PascalCase()
	if !response.Body.IsEmpty() {
		w.Line(`class %s(var body: %s) : %s`, serviceResponseImplementationName, types.Kotlin(&response.Body.Type.Definition), responseInterfaceName(response.Operation))
	} else {
		w.Line(`class %s : %s`, serviceResponseImplementationName, responseInterfaceName(response.Operation))
	}
}

func responseInterfaceName(operation *spec.NamedOperation) string {
	return fmt.Sprintf(`%sResponse`, operation.Name.PascalCase())
}
