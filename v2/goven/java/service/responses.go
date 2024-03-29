package service

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/goven/generator"
	"github.com/specgen-io/specgen/v2/goven/java/types"
	"github.com/specgen-io/specgen/v2/goven/java/writer"
	"github.com/specgen-io/specgen/v2/goven/spec"
)

func (g *Generator) responseInterface(operation *spec.NamedOperation) *generator.CodeFile {
	w := writer.New(g.Packages.ServicesApi(operation.InApi), responseInterfaceName(operation))
	w.Line(`import %s;`, g.Packages.Models(operation.InApi.InHttp.InVersion).PackageStar)
	w.Line(`import %s;`, g.Packages.ErrorsModels.PackageStar)
	w.EmptyLine()
	w.Line(`public interface [[.ClassName]] {`)
	for index, response := range operation.Responses {
		if index > 0 {
			w.EmptyLine()
		}
		responseImpl(w.Indented(), g.Types, &response)
	}
	w.Line(`}`)
	return w.ToCodeFile()
}

func responseImpl(w *writer.Writer, types *types.Types, response *spec.OperationResponse) {
	serviceResponseImplementationName := response.Name.PascalCase()
	w.Line(`class %s implements %s {`, serviceResponseImplementationName, responseInterfaceName(response.Operation))
	if !response.Body.IsEmpty() {
		w.Line(`  public %s body;`, types.Java(&response.Body.Type.Definition))
		w.EmptyLine()
		w.Line(`  public %s() {`, serviceResponseImplementationName)
		w.Line(`  }`)
		w.EmptyLine()
		w.Line(`  public %s(%s body) {`, serviceResponseImplementationName, types.Java(&response.Body.Type.Definition))
		w.Line(`    this.body = body;`)
		w.Line(`  }`)
	}
	w.Line(`}`)
}

func responseInterfaceName(operation *spec.NamedOperation) string {
	return fmt.Sprintf(`%sResponse`, operation.Name.PascalCase())
}

func getResponseBody(response *spec.OperationResponse, responseVarName string) string {
	return fmt.Sprintf(`((%s.%s) %s).body`, responseInterfaceName(response.Operation), response.Name.PascalCase(), responseVarName)
}
