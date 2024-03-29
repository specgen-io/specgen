package client

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/goven/typescript/writer"

	"github.com/specgen-io/specgen/v2/goven/spec"
	"github.com/specgen-io/specgen/v2/goven/typescript/types"
)

func generateOperationResponse(w *writer.Writer, operation *spec.NamedOperation) {
	w.Line("export type %s =", responseTypeName(operation))
	for _, response := range operation.Responses {
		if !response.Body.IsEmpty() {
			w.Line(`  | { status: "%s", data: %s }`, response.Name.Source, types.TsType(&response.Body.Type.Definition))
		} else {
			w.Line(`  | { status: "%s" }`, response.Name.Source)
		}
	}
}

func responseType(operation *spec.NamedOperation) string {
	successResponses := operation.Responses.Success()
	if len(successResponses) == 1 {
		if successResponses[0].Body.IsEmpty() {
			return "void"
		} else {
			return types.TsType(&successResponses[0].Body.Type.Definition)
		}
	} else {
		return responseTypeName(operation)
	}
}

func newResponse(response *spec.Response, body string) string {
	if body == `` {
		return fmt.Sprintf(`{ status: "%s" }`, response.Name.Source)
	} else {
		return fmt.Sprintf(`{ status: "%s", data: %s }`, response.Name.Source, body)
	}
}

func responseTypeName(operation *spec.NamedOperation) string {
	return operation.Name.PascalCase() + "Response"
}
