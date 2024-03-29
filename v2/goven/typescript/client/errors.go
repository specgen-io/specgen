package client

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/goven/generator"
	"github.com/specgen-io/specgen/v2/goven/spec"
	"github.com/specgen-io/specgen/v2/goven/typescript/types"
	"github.com/specgen-io/specgen/v2/goven/typescript/writer"
)

func errorExceptionName(response *spec.Response) string {
	return fmt.Sprintf(`%sException`, response.Name.PascalCase())
}

func (g *CommonGenerator) Errors(httpErrors *spec.HttpErrors) *generator.CodeFile {
	w := writer.New(g.Modules.Errors)
	w.Imports.Star(g.Modules.ErrorsModels, "errors")
	w.Line(`export * from '%s'`, g.Modules.ErrorsModels.GetImport(g.Modules.Errors))
	w.EmptyLine()
	w.Lines(`
export class ResponseException extends Error {
  constructor(message: string) {
    super(message)
  }
}
`)
	for _, response := range httpErrors.Responses {
		w.EmptyLine()
		w.Line(`export class %s extends ResponseException {`, errorExceptionName(&response.Response))
		if response.Body.Is(spec.ResponseBodyEmpty) {
			w.Line(`  constructor() {`)
			w.Line(`    super('Error response with status code %s')`, spec.HttpStatusCode(response.Name))
			w.Line(`  }`)
		} else {
			w.Line(`  public body: %s`, types.TsType(&response.Body.Type.Definition))
			w.EmptyLine()
			w.Line(`  constructor(body: %s) {`, types.TsType(&response.Body.Type.Definition))
			w.Line(`    super('Error response with status code %s')`, spec.HttpStatusCode(response.Name))
			w.Line(`    this.body = body`)
			w.Line(`  }`)
		}
		w.Line(`}`)
	}
	return w.ToCodeFile()
}
