package client

import (
	"github.com/specgen-io/specgen/v2/goven/generator"
	"github.com/specgen-io/specgen/v2/goven/kotlin/packages"
	"github.com/specgen-io/specgen/v2/goven/kotlin/writer"
)

func staticConfigFiles(thePackage, jsonPackage packages.Package) []generator.CodeFile {
	files := []generator.CodeFile{}
	files = append(files, *objectMapperConfig(thePackage, jsonPackage))
	files = append(files, *clientConfig(thePackage))

	return files
}

func objectMapperConfig(thePackage packages.Package, jsonPackage packages.Package) *generator.CodeFile {
	w := writer.New(thePackage, `ObjectMapperConfig`)
	w.Template(
		map[string]string{
			`JsonPackageName`: jsonPackage.PackageName,
		}, `
import com.fasterxml.jackson.databind.ObjectMapper
import com.fasterxml.jackson.module.kotlin.jacksonObjectMapper
import io.micronaut.context.annotation.*
import io.micronaut.jackson.ObjectMapperFactory
import [[.JsonPackageName]].*

@Factory
@Replaces(ObjectMapperFactory::class)
class ObjectMapperConfig {
	@Bean
	@Replaces(ObjectMapper::class)
	fun objectMapper(): ObjectMapper {
		val objectMapper = jacksonObjectMapper()
		setupObjectMapper(objectMapper)
		return objectMapper
	}

	@Bean
	fun json(): Json {
		return Json(objectMapper())
	}
}
`)
	return w.ToCodeFile()
}

func clientConfig(thePackage packages.Package) *generator.CodeFile {
	w := writer.New(thePackage, `ClientConfiguration`)
	w.Lines(`
class ClientConfiguration {
	companion object {
		const val BASE_URL = "http://localhost:8081"
	}
}
`)
	return w.ToCodeFile()
}
