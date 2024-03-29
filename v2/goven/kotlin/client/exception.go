package client

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/goven/generator"
	"github.com/specgen-io/specgen/v2/goven/kotlin/packages"
	"github.com/specgen-io/specgen/v2/goven/kotlin/writer"
	"github.com/specgen-io/specgen/v2/goven/spec"
)

func clientException(thePackage packages.Package) *generator.CodeFile {
	w := writer.New(thePackage, `ClientException`)
	w.Lines(`
import java.lang.RuntimeException

open class [[.ClassName]] : RuntimeException {
	constructor() : super()
	constructor(message: String) : super(message)
	constructor(cause: Throwable) : super(cause)
	constructor(message: String, cause: Throwable) : super(message, cause)
}
`)
	return w.ToCodeFile()
}

func responseException(thePackage packages.Package) *generator.CodeFile {
	w := writer.New(thePackage, `ResponseException`)
	w.Lines(`
import java.lang.RuntimeException

open class [[.ClassName]](message: String) : RuntimeException(message)
`)
	return w.ToCodeFile()
}

func errorResponseException(thePackage, errorsModelsPackage packages.Package, error *spec.Response) *generator.CodeFile {
	w := writer.New(thePackage, errorExceptionClassName(error))
	w.Imports.PackageStar(errorsModelsPackage)
	errorBody := ""
	if !error.Body.Is(spec.ResponseBodyEmpty) {
		errorBody = fmt.Sprintf(`(val body: %s)`, error.Body.Type.Definition)
	}
	w.Line(`class [[.ClassName]]%s : ResponseException("Error response with status code %s")`, errorBody, spec.HttpStatusCode(error.Name))

	return w.ToCodeFile()
}

func errorExceptionClassName(error *spec.Response) string {
	return fmt.Sprintf(`%sException`, error.Name.PascalCase())
}
