package client

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/goven/generator"
	"github.com/specgen-io/specgen/v2/goven/kotlin/packages"
	"github.com/specgen-io/specgen/v2/goven/kotlin/types"
	"github.com/specgen-io/specgen/v2/goven/kotlin/writer"
	"github.com/specgen-io/specgen/v2/goven/spec"
)

func clientException(thePackage packages.Package) *generator.CodeFile {
	w := writer.New(thePackage, `ClientException`)
	w.Lines(`
import java.lang.RuntimeException

open class ClientException : RuntimeException {
	constructor() : super()
	constructor(message: String) : super(message)
	constructor(cause: Throwable) : super(cause)
	constructor(message: String, cause: Throwable) : super(message, cause)
}
`)
	return w.ToCodeFile()
}

func inheritedClientException(thePackage, errorsModelsPackage packages.Package, types *types.Types, error *spec.Response) *generator.CodeFile {
	errorName := types.Kotlin(&error.Type.Definition)
	className := fmt.Sprintf(`%sException`, errorName)
	w := writer.New(thePackage, className)
	w.Template(
		map[string]string{
			`ErrorsModelsPackage`: errorsModelsPackage.PackageName,
			`ErrorName`:           errorName,
		}, `
import [[.ErrorsModelsPackage]].*

class [[.ClassName]](error: [[.ErrorName]]) : ClientException("Body: $error")
`)
	return w.ToCodeFile()
}