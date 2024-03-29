package models

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/goven/generator"
	"github.com/specgen-io/specgen/v2/goven/kotlin/packages"
	"github.com/specgen-io/specgen/v2/goven/kotlin/types"
	"github.com/specgen-io/specgen/v2/goven/kotlin/writer"
	"github.com/specgen-io/specgen/v2/goven/spec"
)

var Moshi = "moshi"

type MoshiGenerator struct {
	generatedSetupMoshiMethods []string
	Types                      *types.Types
	Packages                   *Packages
}

func NewMoshiGenerator(types *types.Types, packages *Packages) *MoshiGenerator {
	return &MoshiGenerator{[]string{}, types, packages}
}

func (g *MoshiGenerator) Models(version *spec.Version) []generator.CodeFile {
	return g.models(version.ResolvedModels, g.Packages.Models(version))
}

func (g *MoshiGenerator) ErrorModels(httperrors *spec.HttpErrors) []generator.CodeFile {
	return g.models(httperrors.ResolvedModels, g.Packages.ErrorsModels)
}

func (g *MoshiGenerator) models(models []*spec.NamedModel, modelsPackage packages.Package) []generator.CodeFile {
	w := writer.New(modelsPackage, `models`)
	w.Imports.Add(g.modelsDefinitionsImports()...)
	w.Imports.Add(g.Types.Imports()...)

	for _, model := range models {
		w.EmptyLine()
		if model.IsObject() {
			g.modelObject(w, model)
		} else if model.IsOneOf() {
			g.modelOneOf(w, model)
		} else if model.IsEnum() {
			g.modelEnum(w, model)
		}
	}

	files := []generator.CodeFile{}
	files = append(files, *w.ToCodeFile())

	g.generatedSetupMoshiMethods = append(g.generatedSetupMoshiMethods, fmt.Sprintf(`%s.setupModelsMoshiAdapters`, modelsPackage.PackageName))
	for range g.generatedSetupMoshiMethods {
		files = append(files, *g.setupOneOfAdapters(models, modelsPackage))
	}

	return files
}

func (g *MoshiGenerator) modelObject(w *writer.Writer, model *spec.NamedModel) {
	className := model.Name.PascalCase()
	w.Line(`data class %s(`, className)
	for _, field := range model.Object.Fields {
		w.Line(`  @Json(name = "%s")`, field.Name.Source)
		w.Line(`  val %s: %s,`, field.Name.CamelCase(), g.Types.Kotlin(&field.Type.Definition))
	}
	w.Line(`)`)
}

func (g *MoshiGenerator) modelEnum(w *writer.Writer, model *spec.NamedModel) {
	enumName := model.Name.PascalCase()
	w.Line(`enum class %s {`, enumName)
	for _, enumItem := range model.Enum.Items {
		w.Line(`  @Json(name = "%s")`, enumItem.Value)
		w.Line(`  %s,`, enumItem.Name.UpperCase())
	}
	w.Line(`}`)
}

func (g *MoshiGenerator) modelOneOf(w *writer.Writer, model *spec.NamedModel) {
	sealedClassName := model.Name.PascalCase()
	w.Line(`sealed class %s {`, sealedClassName)
	for _, item := range model.OneOf.Items {
		w.Line(`  data class %s(val data: %s): %s()`, oneOfItemClassName(&item), g.Types.Kotlin(&item.Type.Definition), sealedClassName)
	}
	w.Line(`}`)
}

func (g *MoshiGenerator) ReadJson(varJson string, typ *spec.TypeDef) string {
	switch typ.Node {
	case spec.ArrayType:
		return fmt.Sprintf(`readList(%s)`, varJson)
	case spec.MapType:
		return fmt.Sprintf(`readMap(%s)`, varJson)
	default:
		return fmt.Sprintf(`read(%s)`, varJson)
	}
}

func (g *MoshiGenerator) WriteJson(varData string, typ *spec.TypeDef) string {
	switch typ.Node {
	case spec.ArrayType:
		return fmt.Sprintf(`writeList(%s)`, varData)
	case spec.MapType:
		return fmt.Sprintf(`writeMap(%s)`, varData)
	default:
		return fmt.Sprintf(`write(%s)`, varData)
	}
}

func (g *MoshiGenerator) modelsDefinitionsImports() []string {
	return []string{
		`com.squareup.moshi.*`,
	}
}

func (g *MoshiGenerator) ModelsUsageImports() []string {
	return []string{
		`com.squareup.moshi.Moshi`,
		`com.squareup.moshi.Types`,
		`com.squareup.moshi.adapter`,
		`com.squareup.moshi.JsonDataException`,
		`java.lang.reflect.ParameterizedType`,
	}
}

func (g *MoshiGenerator) ValidationErrorsHelpers() *generator.CodeFile {
	w := writer.New(g.Packages.Errors, `ValidationErrorsHelpers`)
	w.Template(
		map[string]string{
			`JsonPackage`:         g.Packages.Json.PackageName,
			`ErrorsModelsPackage`: g.Packages.ErrorsModels.PackageName,
		}, `
import [[.JsonPackage]].*
import [[.ErrorsModelsPackage]].*
import java.util.regex.Pattern

object ValidationErrorsHelpers {
	private val pathPattern = Pattern.compile("\\$\\.([^ ]+)")

    fun extractValidationErrors(exception: JsonParseException): List<ValidationError>? {
        val matcher = pathPattern.matcher(exception.message!!)
        if (matcher.find()) {
            val jsonPath = matcher.group(1)
            return listOf(ValidationError(jsonPath, "parsing_failed", exception.message))
        }
        return null
    }
}
`)
	return w.ToCodeFile()
}

func (g *MoshiGenerator) JsonHelpers() []generator.CodeFile {
	files := []generator.CodeFile{}

	files = append(files, *g.json())
	files = append(files, *g.jsonParseException())
	files = append(files, g.setupLibrary()...)

	return files
}

func (g *MoshiGenerator) json() *generator.CodeFile {
	w := writer.New(g.Packages.Json, `Json`)
	w.Lines(`
import com.squareup.moshi.*
import java.lang.reflect.ParameterizedType
import java.io.Reader

class Json(private val moshi: Moshi) {
	fun <T> write(type: Class<T>, data: T): String {
		return moshi.adapter(type).toJson(data)
	}

	inline fun <reified T> write(data: T) : String {
		return write(T::class.java, data)
	}

	fun <T> write(type: ParameterizedType, data: T): String {
		return moshi.adapter<Any>(type).toJson(data)
	}

	inline fun <reified T> writeList(data: List<T>): String {
		return write(Types.newParameterizedType(List::class.java, T::class.java), data)
	}

	inline fun <reified K, reified V> writeMap(data: Map<K, V>): String {
		return write(Types.newParameterizedType(MutableMap::class.java, K::class.java, V::class.java), data)
	}

	fun <T> read(jsonStr: String, type: Class<T>): T {
		return try {
			moshi.adapter(type).fromJson(jsonStr)!!
		} catch (exception: JsonDataException) {
			throw JsonParseException(exception)
		}
	}

	fun <T> read(jsonStr: String, type: ParameterizedType): T {
		return try {
			moshi.adapter<T>(type).fromJson(jsonStr)!!
		} catch (exception: JsonDataException) {
			throw JsonParseException(exception)
		}
	}

	fun <T> read(reader: Reader, type: Class<T>): T {
		return try {
			moshi.adapter(type).fromJson(reader.readText())!!
		} catch (exception: JsonDataException) {
			throw JsonParseException(exception)
		}
	}

	fun <T> read(reader: Reader, type: ParameterizedType): T {
		return try {
			moshi.adapter<T>(type).fromJson(reader.readText())!!
		} catch (exception: JsonDataException) {
			throw JsonParseException(exception)
		}
	}

	inline fun <reified T> read(jsonStr: String): T {
		return read(jsonStr, T::class.java)
	}

	inline fun <reified T> readList(jsonStr: String): List<T> {
		return read(jsonStr, Types.newParameterizedType(List::class.java, T::class.java))
	}

	inline fun <reified K, reified V> readMap(jsonStr: String): Map<K, V> {
		return read(jsonStr, Types.newParameterizedType(MutableMap::class.java, K::class.java, V::class.java))
	}

	inline fun <reified T> read(reader: Reader): T {
		return read(reader, T::class.java)
	}

	inline fun <reified T> readList(reader: Reader): List<T> {
		return read(reader, Types.newParameterizedType(List::class.java, T::class.java))
	}

	inline fun <reified K, reified V> readMap(reader: Reader): Map<K, V> {
		return read(reader, Types.newParameterizedType(MutableMap::class.java, K::class.java, V::class.java))
	}
}
`)
	return w.ToCodeFile()
}

func (g *MoshiGenerator) jsonParseException() *generator.CodeFile {
	w := writer.New(g.Packages.Json, `JsonParseException`)
	w.Lines(`
class JsonParseException(exception: Throwable) :
	RuntimeException("Failed to parse body: " + exception.message, exception)
`)
	return w.ToCodeFile()
}

func (g *MoshiGenerator) setupLibrary() []generator.CodeFile {
	files := []generator.CodeFile{}
	files = append(files, *g.setupAdapters())
	files = append(files, *bigDecimalAdapter(g.Packages.JsonAdapters))
	files = append(files, *localDateAdapter(g.Packages.JsonAdapters))
	files = append(files, *localDateTimeAdapter(g.Packages.JsonAdapters))
	files = append(files, *uuidAdapter(g.Packages.JsonAdapters))
	files = append(files, *unionAdapterFactory(g.Packages.JsonAdapters))
	files = append(files, *unwrapFieldAdapterFactory(g.Packages.JsonAdapters))
	return files
}

func (g *MoshiGenerator) setupAdapters() *generator.CodeFile {
	w := writer.New(g.Packages.Json, `CustomMoshiAdapters`)
	w.Imports.Add(`com.squareup.moshi.Moshi`)
	w.Imports.Add(`com.squareup.moshi.kotlin.reflect.KotlinJsonAdapterFactory`)
	w.Imports.PackageStar(g.Packages.JsonAdapters)
	w.Lines(`
fun setupMoshiAdapters(moshiBuilder: Moshi.Builder): Moshi.Builder {
	moshiBuilder
		.add(BigDecimalAdapter())
		.add(UuidAdapter())
		.add(LocalDateAdapter())
		.add(LocalDateTimeAdapter())
`)
	w.EmptyLine()
	for _, setupMoshiMethod := range g.generatedSetupMoshiMethods {
		w.Line(`	%s(moshiBuilder);`, setupMoshiMethod)
	}
	w.EmptyLine()
	w.Lines(`
	moshiBuilder
		.add(KotlinJsonAdapterFactory())

	return moshiBuilder;
}
`)
	return w.ToCodeFile()
}

func (g *MoshiGenerator) setupOneOfAdapters(models []*spec.NamedModel, modelsPackage packages.Package) *generator.CodeFile {
	w := writer.New(modelsPackage, `ModelsMoshiAdapters`)
	w.Imports.Add(`com.squareup.moshi.Moshi`)
	w.Imports.PackageStar(g.Packages.JsonAdapters)
	w.EmptyLine()
	w.Line(`fun setupModelsMoshiAdapters(moshiBuilder: Moshi.Builder) {`)
	for _, model := range models {
		if model.IsOneOf() {
			w.Indent()
			w.Line(`moshiBuilder`)
			modelName := model.Name.PascalCase()
			for _, item := range model.OneOf.Items {
				w.Line(`  .add(UnwrapFieldAdapterFactory(%s.%s::class.java))`, modelName, oneOfItemClassName(&item))
			}
			addUnionAdapterFactory := fmt.Sprintf(`  .add(UnionAdapterFactory.of(%s::class.java)`, modelName)
			if model.OneOf.Discriminator != nil {
				w.Line(`%s.withDiscriminator("%s")`, addUnionAdapterFactory, *model.OneOf.Discriminator)
			} else {
				w.Line(addUnionAdapterFactory)
			}
			for _, item := range model.OneOf.Items {
				w.Line(`    .withSubtype(%s.%s::class.java, "%s")`, modelName, oneOfItemClassName(&item), item.Name.Source)
			}
			w.Line(`)`)
			w.Unindent()
		}
	}
	w.Line(`}`)
	return w.ToCodeFile()
}

func bigDecimalAdapter(thePackage packages.Package) *generator.CodeFile {
	w := writer.New(thePackage, `BigDecimalAdapter`)
	w.Line(`
import com.squareup.moshi.*
import okio.*

import java.io.ByteArrayInputStream
import java.math.BigDecimal
import java.nio.charset.StandardCharsets

class BigDecimalAdapter {
	@FromJson
	fun fromJson(reader: JsonReader): BigDecimal {
		val token = reader.peek()
		if (token != JsonReader.Token.NUMBER) {
			throw JsonDataException("BigDecimal should be represented as number in JSON, found: " + token.name)
		}
		val source = reader.nextSource()
		return BigDecimal(String(source.readByteArray(), StandardCharsets.UTF_8))
	}

	@ToJson
	fun toJson(writer: JsonWriter, value: BigDecimal) {
		val source = ByteArrayInputStream(value.toString().toByteArray()).source()
		val buffer = source.buffer()
		writer.value(buffer)
	}
}
`)
	return w.ToCodeFile()
}

func localDateAdapter(thePackage packages.Package) *generator.CodeFile {
	w := writer.New(thePackage, `LocalDateAdapter`)
	w.Lines(`
import com.squareup.moshi.*
import java.time.*

class LocalDateAdapter {
	@FromJson
	private fun fromJson(string: String): LocalDate {
		return LocalDate.parse(string)
	}

	@ToJson
	private fun toJson(value: LocalDate): String {
		return value.toString()
	}
}
`)
	return w.ToCodeFile()
}

func localDateTimeAdapter(thePackage packages.Package) *generator.CodeFile {
	w := writer.New(thePackage, `LocalDateTimeAdapter`)
	w.Lines(`
import com.squareup.moshi.*

import java.time.*

class LocalDateTimeAdapter {
	@FromJson
	private fun fromJson(string: String): LocalDateTime {
		return LocalDateTime.parse(string)
	}

    @ToJson
	private fun toJson(value: LocalDateTime): String {
		return value.toString()
	}
}
`)
	return w.ToCodeFile()
}

func uuidAdapter(thePackage packages.Package) *generator.CodeFile {
	w := writer.New(thePackage, `UuidAdapter`)
	w.Lines(`
import com.squareup.moshi.*
import java.util.UUID

class UuidAdapter {
	@FromJson
	private fun fromJson(string: String): UUID {
		return UUID.fromString(string)
	}

	@ToJson
	private fun toJson(value: UUID): String {
		return value.toString()
	}
}
`)
	return w.ToCodeFile()
}

func unionAdapterFactory(thePackage packages.Package) *generator.CodeFile {
	w := writer.New(thePackage, `UnionAdapterFactory`)
	w.Lines(`
import com.squareup.moshi.*
import java.lang.reflect.Type

class UnionAdapterFactory<T> internal constructor(
	private val baseType: Class<T>,
	private val discriminator: String?,
	private val tags: List<String>,
	private val subtypes: List<Type>,
	private val fallbackAdapter: JsonAdapter<Any>?
) : JsonAdapter.Factory {

	companion object {
		fun <T> of(baseType: Class<T>): UnionAdapterFactory<T> {
			return UnionAdapterFactory(baseType, null, emptyList(), emptyList(), null)
		}
	}

	fun withDiscriminator(discriminator: String?): UnionAdapterFactory<T> {
		if (discriminator == null) throw NullPointerException("discriminator == null")
		return UnionAdapterFactory(baseType, discriminator, tags, subtypes, fallbackAdapter)
	}

	fun withSubtype(subtype: Class<out T>, tag: String): UnionAdapterFactory<T> {
		require(!tags.contains(tag)) { "Tags must be unique." }
		val newTags: MutableList<String> = ArrayList(tags).also {
			it.add(tag)
		}
		val newSubtypes: MutableList<Type> = ArrayList(subtypes).also {
			it.add(subtype)
		}
		return UnionAdapterFactory(baseType, discriminator, newTags, newSubtypes, fallbackAdapter)
	}

	private fun withFallbackAdapter(fallbackJsonAdapter: JsonAdapter<Any>?): UnionAdapterFactory<T> {
		return UnionAdapterFactory(baseType, discriminator, tags, subtypes, fallbackJsonAdapter)
	}

	fun withDefaultValue(defaultValue: T): UnionAdapterFactory<T> {
		return withFallbackAdapter(buildFallbackAdapter(defaultValue))
	}

	private fun buildFallbackAdapter(defaultValue: T): JsonAdapter<Any> {
		return object : JsonAdapter<Any>() {
			override fun fromJson(reader: JsonReader): Any? {
				reader.skipValue()
				return defaultValue
			}

			override fun toJson(writer: JsonWriter, value: Any?) {
				throw IllegalArgumentException("Expected one of " + subtypes + " but found " + value + ", a " + value!!.javaClass + ". Register this subtype.")
			}
		}
	}

	override fun create(type: Type, annotations: Set<Annotation>, moshi: Moshi): JsonAdapter<*>? {
		if (Types.getRawType(type) != baseType || annotations.isNotEmpty()) {
			return null
		}

		val jsonAdapters: MutableList<JsonAdapter<Any>> = java.util.ArrayList(subtypes.size)

		for (element in subtypes) {
			jsonAdapters.add(moshi.adapter(element))
		}

		return if (discriminator != null) {
			UnionDiscriminatorAdapter(discriminator, tags, subtypes, jsonAdapters, fallbackAdapter).nullSafe()
		} else {
			UnionWrapperAdapter(tags, subtypes, jsonAdapters, fallbackAdapter).nullSafe()
		}
	}

	internal class UnionDiscriminatorAdapter(
		private val discriminator: String,
		private val tags: List<String>,
		private val subtypes: List<Type>,
		private val adapters: List<JsonAdapter<Any>>,
		private val fallbackAdapter: JsonAdapter<Any>?
    ) : JsonAdapter<Any>() {
		private val discriminatorOptions: JsonReader.Options = JsonReader.Options.of(discriminator)
		private val tagsOptions: JsonReader.Options = JsonReader.Options.of(*tags.toTypedArray())

		override fun fromJson(reader: JsonReader): Any? {
			val tagIndex = getTagIndex(reader)
			var adapter = fallbackAdapter
			if (tagIndex != -1) {
				adapter = adapters[tagIndex]
			}
			return adapter!!.fromJson(reader)
		}

		override fun toJson(writer: JsonWriter, value: Any?) {
			val tagIndex = getTagIndex(value!!)
			if (tagIndex == -1) {
				fallbackAdapter!!.toJson(writer, value)
			} else {
				val adapter = adapters[tagIndex]
				writer.beginObject()
				writer.name(discriminator).value(tags[tagIndex])
				val flattenToken = writer.beginFlatten()
				adapter.toJson(writer, value)
				writer.endFlatten(flattenToken)
				writer.endObject()
			}
		}

		private fun getTagIndex(reader: JsonReader): Int {
			val peeked = reader.peekJson()
			peeked.setFailOnUnknown(false)
			peeked.use {
				it.beginObject()
				while (it.hasNext()) {
					if (it.selectName(discriminatorOptions) == -1) {
						it.skipName()
						it.skipValue()
						continue
				}
					val tagIndex = it.selectString(tagsOptions)
					if (tagIndex == -1 && fallbackAdapter == null) {
						throw JsonDataException("Expected one of " + tags + " for key '" + discriminator + "' but found '" + it.nextString() + "'. Register a subtype for this tag.")
					}
					return tagIndex
				}
				throw JsonDataException("Missing discriminator field $discriminator")
			}
		}

		private fun getTagIndex(value: Any): Int {
			val type: Class<*> = value.javaClass
			val tagIndex = subtypes.indexOf(type)
			return if (tagIndex == -1 && fallbackAdapter == null) {
				throw IllegalArgumentException("Expected one of " + subtypes + " but found " + value + ", a " + value.javaClass + ". Register this subtype.")
			} else {
				tagIndex
			}
		}

		override fun toString(): String {
			return "UnionDiscriminatorAdapter($discriminator)"
		}
	}

	internal class UnionWrapperAdapter(
		private val tags: List<String>,
		private val subtypes: List<Type>,
		private val adapters: List<JsonAdapter<Any>>,
		private val fallbackAdapter: JsonAdapter<Any>?
	) : JsonAdapter<Any>() {
		private val tagsOptions: JsonReader.Options = JsonReader.Options.of(*tags.toTypedArray())

		override fun fromJson(reader: JsonReader): Any? {
			val tagIndex: Int = getTagIndex(reader)
			return if (tagIndex == -1) {
				fallbackAdapter!!.fromJson(reader)
			} else {
				reader.beginObject()
				reader.skipName()
				val value = adapters[tagIndex].fromJson(reader)
				reader.endObject()
				value
			}
		}

		override fun toJson(writer: JsonWriter, value: Any?) {
			val tagIndex: Int = getTagIndex(value!!)
			if (tagIndex == -1) {
				fallbackAdapter!!.toJson(writer, value)
			} else {
				val adapter = adapters[tagIndex]
				writer.beginObject()
				writer.name(tags[tagIndex])
				adapter.toJson(writer, value)
				writer.endObject()
			}
		}

		private fun getTagIndex(reader: JsonReader): Int {
			val peeked = reader.peekJson()
			peeked.setFailOnUnknown(false)
			return peeked.use {
				it.beginObject()
				val tagIndex = it.selectName(tagsOptions)
				if (tagIndex == -1 && fallbackAdapter == null) {
					throw JsonDataException("Expected one of keys:" + tags + "' but found '" + it.nextString() + "'. Register a subtype for this tag.")
				}
				tagIndex
			}
		}

		private fun getTagIndex(value: Any): Int {
			val type: Class<*> = value.javaClass
			val tagIndex = subtypes.indexOf(type)
			require(!(tagIndex == -1 && fallbackAdapter == null)) { "Expected one of " + subtypes + " but found " + value + ", a " + value.javaClass + ". Register this subtype." }
			return tagIndex
		}

		override fun toString(): String {
			return "UnionWrapperAdapter"
		}
	}
}
`)
	return w.ToCodeFile()
}

func unwrapFieldAdapterFactory(thePackage packages.Package) *generator.CodeFile {
	w := writer.New(thePackage, `UnwrapFieldAdapterFactory`)
	w.Lines(`
import com.squareup.moshi.*
import java.io.IOException
import java.lang.reflect.*

class UnwrapFieldAdapterFactory<T>(private val type: Class<T>) : JsonAdapter.Factory {
	override fun create(type: Type, annotations: Set<Annotation?>, moshi: Moshi): JsonAdapter<*>? {
		if (Types.getRawType(type) != this.type || annotations.isNotEmpty()) {
			return null
		}

		val fields = this.type.declaredFields
		if (fields.size != 1) {
			throw RuntimeException("Type " + type.typeName + " has " + fields.size + " fields, unwrap adapter can be used only with single-field types")
		}
		val field = fields[0]
		val getterName = "get" + field.name.replaceFirstChar { it.uppercase() }

		val getter = try {
			this.type.getDeclaredMethod(getterName)
		} catch (e: NoSuchMethodException) {
			throw RuntimeException("Type " + type.typeName + " field " + field.name + " does not have getter method " + field.type.name + ", it's required for unwrap adapter", e)
		}

		val constructor: Constructor<T> = try {
			this.type.getDeclaredConstructor(field.type)
		} catch (e: NoSuchMethodException) {
			throw RuntimeException("Type " + type.typeName + " does not have constructor with single parameter of type " + field.type.name + ", it's required for unwrap adapter")
		}

		val fieldAdapter: JsonAdapter<*>  = moshi.adapter(field.type)
		return UnwrapFieldAdapter(constructor, getter, fieldAdapter)
	}

	inner class UnwrapFieldAdapter<O, I>(
		private val constructor: Constructor<O>,
		private val getter: Method,
		private val fieldAdapter: JsonAdapter<I>
	) : JsonAdapter<Any>() {

		override fun fromJson(reader: JsonReader): Any? {
			val fieldValue = fieldAdapter.fromJson(reader)
			return try {
				constructor.newInstance(fieldValue)
			} catch (e: Throwable) {
				throw IOException("Failed to create object with constructor " + constructor.name, e)
			}
		}

		@Suppress("UNCHECKED_CAST")
		override fun toJson(writer: JsonWriter, value: Any?) {
			if (value != null) {
				val fieldValue: I = try {
					getter.invoke(value) as I
				} catch (e: IllegalAccessException) {
					throw IOException("Failed to get value of field " + getter.name, e)
				}
				fieldAdapter.toJson(writer, fieldValue)
			} else {
				fieldAdapter.toJson(writer, null)
			}
		}

		override fun toString(): String {
			return "UnwrapFieldAdapter(" + getter.name + ")"
		}
	}
}
`)
	return w.ToCodeFile()
}

func (g *MoshiGenerator) JsonMapperInit() string {
	return fmt.Sprintf(`setupMoshiAdapters(Moshi.Builder()).build()`)
}

func (g *MoshiGenerator) JsonMapperType() string {
	return `Moshi`
}
