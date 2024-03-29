package generators

import (
	"github.com/specgen-io/specgen/v2/goven/scala/writer"
	"github.com/specgen-io/specgen/v2/goven/spec"
)

func responseType(operation *spec.NamedOperation) string {
	if len(operation.Responses) == 1 {
		response := operation.Responses[0]
		return ResponseBodyScalaType(&response.Body)
	} else {
		return responseTypeName(operation)
	}
}

func responseTypeName(operation *spec.NamedOperation) string {
	return operation.Name.PascalCase() + "Response"
}

func generateResponse(w *writer.Writer, operation *spec.NamedOperation) {
	if len(operation.Responses) > 1 {
		w.Line(`sealed trait %s`, responseTypeName(operation))
		w.Line(`object %s {`, responseTypeName(operation))
		for _, response := range operation.Responses {
			var bodyParam = ""
			if !response.Body.IsEmpty() {
				bodyParam = `body: ` + ScalaType(&response.Body.Type.Definition)
			}
			w.Line(`  case class %s(%s) extends %s`, response.Name.PascalCase(), bodyParam, responseTypeName(operation))
		}
		w.Line(`}`)
	}
}
