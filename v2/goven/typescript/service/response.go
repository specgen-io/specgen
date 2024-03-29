package service

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/goven/typescript/writer"

	"github.com/specgen-io/specgen/v2/goven/spec"
	"github.com/specgen-io/specgen/v2/goven/typescript/types"
)

func GenerateOperationResponse(w *writer.Writer, operation *spec.NamedOperation) {
	w.Line("export type %s =", responseTypeName(operation))
	for _, response := range operation.Responses {
		if !response.Body.IsEmpty() {
			w.Line(`  | { status: "%s", data: %s }`, response.Name.Source, types.TsType(&response.Body.Type.Definition))
		} else {
			w.Line(`  | { status: "%s" }`, response.Name.Source)
		}
	}
}

func ResponseType(operation *spec.NamedOperation, servicePackage string) string {
	if len(operation.Responses) == 1 {
		response := operation.Responses[0]
		if response.Body.IsEmpty() {
			return "void"
		}
		return types.TsType(&response.Body.Type.Definition)
	}
	result := responseTypeName(operation)
	if servicePackage != "" {
		result = "service." + result
	}
	return result
}

func New(response *spec.Response, body string) string {
	if body == `` {
		return fmt.Sprintf(`{ status: "%s" }`, response.Name.Source)
	} else {
		return fmt.Sprintf(`{ status: "%s", data: %s }`, response.Name.Source, body)
	}
}

func responseTypeName(operation *spec.NamedOperation) string {
	return operation.Name.PascalCase() + "Response"
}
