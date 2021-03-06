package genscala

import (
	"fmt"
	spec "github.com/specgen-io/spec"
	"github.com/specgen-io/specgen/v2/gen"
	"github.com/vsapronov/gopoetry/scala"
	"path/filepath"
	"strings"
)

func GenerateSttpClient(serviceFile string, generatePath string) error {
	specification, err := spec.ReadSpec(serviceFile)
	if err != nil {
		return err
	}

	clientPackage := clientPackageName(specification.Name)

	scalaCirceFile := generateJson("spec", filepath.Join(generatePath, "Json.scala"))
	scalaHttpStaticFile := generateStringParams("spec", filepath.Join(generatePath, "StringParams.scala"))

	modelsFiles := GenerateCirceModels(specification, clientPackage, generatePath)
	interfacesFiles := generateClientInterfaces(specification, clientPackage, generatePath)
	implsFiles := generateClientImplementations(specification, clientPackage, generatePath)

	sourceManaged := append(modelsFiles, *scalaCirceFile)
	sourceManaged = append(sourceManaged, *scalaHttpStaticFile)
	sourceManaged = append(sourceManaged, interfacesFiles...)
	sourceManaged = append(sourceManaged, implsFiles...)

	err = gen.WriteFiles(sourceManaged, true)
	if err != nil {
		return err
	}

	return nil
}

func clientPackageName(name spec.Name) string {
	return name.FlatCase() + ".client"
}

func generateClientImplementations(specification *spec.Spec, packageName string, outPath string) []gen.TextFile {
	files := []gen.TextFile{}
	for _, version := range specification.Versions {
		versionFile := generateClientApiImplementations(&version, packageName, outPath)
		files = append(files, *versionFile)
	}
	return files
}

func generateClientApiImplementations(version *spec.Version, packageName string, outPath string) *gen.TextFile {
	unit := Unit(versionedPackage(version.Version, packageName))

	unit.
		Import("scala.concurrent._").
		Import("org.slf4j._").
		Import("com.softwaremill.sttp._").
		Import("spec.Jsoner").
		Import("spec.ParamsTypesBindings._")

	for _, api := range version.Http.Apis {
		apiTrait := generateClientApiClass(api)
		unit.AddDeclarations(apiTrait)
	}

	return &gen.TextFile{
		Path:    filepath.Join(outPath, version.Version.PascalCase()+"Client.scala"),
		Content: unit.Code(),
	}
}

func generateClientInterfaces(specification *spec.Spec, packageName string, outPath string) []gen.TextFile {
	files := []gen.TextFile{}
	for _, version := range specification.Versions {
		versionFile := generateClientApisInterfaces(&version, packageName, outPath)
		files = append(files, *versionFile)
	}
	return files
}

func generateClientApisInterfaces(version *spec.Version, packageName string, outPath string) *gen.TextFile {
	unit := Unit(versionedPackage(version.Version, packageName))

	unit.
		Import("scala.concurrent._")

	for _, api := range version.Http.Apis {
		apiTrait := generateClientApiTrait(api)
		unit.AddDeclarations(apiTrait)
	}

	for _, api := range version.Http.Apis {
		apiObject := generateApiInterfaceResponse(api, clientTraitName(api.Name))
		unit.AddDeclarations(apiObject)
	}

	return &gen.TextFile{
		Path:    filepath.Join(outPath, version.Version.PascalCase()+"Interfaces.scala"),
		Content: unit.Code(),
	}
}

func createParams(params []spec.NamedParam, defaulted bool) []scala.Writable {
	methodParams := []scala.Writable{}
	for _, param := range params {
		if !defaulted && param.Default == nil {
			methodParams = append(methodParams, Param(param.Name.CamelCase(), ScalaType(&param.Type.Definition)))
		}
		if defaulted && param.Default != nil {
			defaultValue := DefaultValue(&param.Type.Definition, *param.Default)
			methodParams = append(methodParams, Param(param.Name.CamelCase(), ScalaType(&param.Type.Definition)).Init(Code(defaultValue)))
		}
	}
	return methodParams
}

func createBodyParam(operation spec.NamedOperation) scala.Writable {
	if operation.Body == nil {
		return nil
	}
	return Param("body", ScalaType(&operation.Body.Type.Definition))
}

func createUrlParams(urlParams []spec.NamedParam) []scala.Writable {
	methodParams := []scala.Writable{}
	for _, param := range urlParams {
		methodParams = append(methodParams, Param(param.Name.CamelCase(), ScalaType(&param.Type.Definition)))
	}
	return methodParams
}

func generateClientOperationSignature(operation spec.NamedOperation) *scala.MethodDeclaration {
	returnType := "Future[" + responseType(operation) + "]"
	method :=
		Def(operation.Name.CamelCase()).Returns(returnType).
			AddParams(createParams(operation.HeaderParams, false)...).
			AddParams(createBodyParam(operation)).
			AddParams(createUrlParams(operation.Endpoint.UrlParams)...).
			AddParams(createParams(operation.QueryParams, false)...).
			AddParams(createParams(operation.HeaderParams, true)...).
			AddParams(createParams(operation.QueryParams, true)...)
	return method
}

func generateClientApiTrait(api spec.Api) *scala.TraitDeclaration {
	apiTraitName := clientTraitName(api.Name)
	apiTrait := Trait(apiTraitName).Add(Import(apiTraitName + "._"))
	for _, operation := range api.Operations {
		apiTrait.Add(generateClientOperationSignature(operation))
	}
	return apiTrait
}

func clientTraitName(apiName spec.Name) string {
	return "I" + apiName.PascalCase() + "Client"
}

func clientClassName(apiName spec.Name) string {
	return apiName.PascalCase() + "Client"
}

func addParamsWriting(params []spec.NamedParam, paramsName string) *scala.StatementsDeclaration {
	code := Statements()
	if params != nil && len(params) > 0 {
		code.Add(Line("val %s = new StringParamsWriter()", paramsName))
		for _, p := range params {
			code.Add(Line(`%s.write("%s", %s)`, paramsName, p.Name.Source, p.Name.CamelCase()))
		}
	}
	return code
}

func generateResponseCases(operation spec.NamedOperation) *scala.StatementsDeclaration {
	cases := scala.Statements()
	for _, response := range operation.Responses {
		responseParam := ``
		if !response.Type.Definition.IsEmpty() {
			responseParam = fmt.Sprintf(`Jsoner.readThrowing[%s](body)`, ScalaType(&response.Type.Definition))
		}
		cases.Add(Line(`case %s => %s.%s(%s)`, spec.HttpStatusCode(response.Name), responseType(operation), response.Name.PascalCase(), responseParam))
	}
	return cases
}

func generateClientOperationImplementation(operation spec.NamedOperation) *scala.StatementsDeclaration {
	httpMethod := strings.ToLower(operation.Endpoint.Method)
	url := operation.FullUrl()
	for _, param := range operation.Endpoint.UrlParams {
		url = strings.Replace(url, spec.UrlParamStr(param.Name.Source), "$"+param.Name.CamelCase(), -1)
	}

	code := Statements(
		addParamsWriting(operation.QueryParams, "query"),
		Statements(Dynamic(func(code *scala.WritableList) {
			if operation.QueryParams != nil && len(operation.QueryParams) > 0 {
				code.Add(Line(`val url = Uri.parse(baseUrl+s"%s").get.params(query.params:_*)`, url))
			} else {
				code.Add(Line(`val url = Uri.parse(baseUrl+s"%s").get`, url))
			}
		})...),
		addParamsWriting(operation.HeaderParams, "headers"),
		Statements(Dynamic(func(code *scala.WritableList) {
			if operation.Body != nil {
				code.Add(
					Line(`val bodyJson = Jsoner.write(body)`),
					Line(`logger.debug(s"Request to url: ${url}, body: ${bodyJson}")`),
				)
			} else {
				code.Add(
					Line(`logger.debug(s"Request to url: ${url}")`),
				)
			}
		})...),
		Line("val response: Future[Response[String]] ="),
		Block(
			Line("sttp"),
			Block(
				Line(`.%s(url)`, httpMethod),
				Statements(Dynamic(func(code *scala.WritableList) {
					if operation.HeaderParams != nil && len(operation.HeaderParams) > 0 {
						code.Add(
							Line(`.headers(headers.params:_*)`),
						)
					}
					if operation.Body != nil {
						code.Add(
							Line(`.header("Content-Type", "application/json")`),
							Line(`.body(bodyJson)`),
						)
					}
				})...),
				Line(`.parseResponseIf { status => status < 500 }`),
				Line(`.send()`),
			),
		),
		Code(`response.map `),
		Scope(
			Line(`response: Response[String] =>`),
			Block(
				Code(`response.body match `),
				Scope(
					Code(`case Right(body) => `),
					Block(
						Line(`logger.debug(s"Response status: ${response.code}, body: ${body}")`),
						Code(`response.code match `),
						Scope(generateResponseCases(operation)),
					),
					Code(`case Left(errorData) => `),
					Block(
						Line(`val errorMessage = s"Request failed, status code: ${response.code}, body: ${new String(errorData)}"`),
						Line(`logger.error(errorMessage)`),
						Line(`throw new RuntimeException(errorMessage)`),
					),
				),
			),
		),
	)
	return code
}

func generateClientApiClass(api spec.Api) *scala.ClassDeclaration {
	apiClassName := clientClassName(api.Name)
	apiTraitName := clientTraitName(api.Name)
	apiClass :=
		Class(apiClassName).Extends(apiTraitName).
			Constructor(Constructor().
				Param("baseUrl", "String").
				ImplicitParam("backend", "SttpBackend[Future, Nothing]"),
			).
			Add(Import(apiTraitName + "._")).
			Add(Import("ExecutionContext.Implicits.global")).
			Add(Line("private val logger: Logger = LoggerFactory.getLogger(this.getClass)"))
	for _, operation := range api.Operations {
		method := generateClientOperationSignature(operation).Body(generateClientOperationImplementation(operation))
		apiClass.Add(method)
	}
	return apiClass
}
