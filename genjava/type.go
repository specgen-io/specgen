package genjava

import (
	"fmt"
	"github.com/specgen-io/spec"
)

func JavaType(typ *spec.TypeDef) string {
	return javaType(typ, false)
}

func javaType(typ *spec.TypeDef, referenceTypesOnly bool) string {
	switch typ.Node {
	case spec.PlainType:
		return PlainJavaType(typ.Plain, referenceTypesOnly)
	case spec.NullableType:
		result := javaType(typ.Child, true)
		return result
	case spec.ArrayType:
		child := javaType(typ.Child, false)
		result := child + "[]"
		return result
	case spec.MapType:
		child := javaType(typ.Child, true)
		result := "Map<String, " + child + ">"
		return result
	default:
		panic(fmt.Sprintf("Unknown type: %v", typ))
	}
}

func PlainJavaType(typ string, referenceTypesOnly bool) string {
	switch typ {
	case spec.TypeInt32:
		if referenceTypesOnly {
			return "Integer"
		} else {
			return "int"
		}
	case spec.TypeInt64:
		if referenceTypesOnly {
			return "Long"
		} else {
			return "long"
		}
	case spec.TypeFloat:
		if referenceTypesOnly {
			return "Float"
		} else {
			return "float"
		}
	case spec.TypeDouble:
		if referenceTypesOnly {
			return "Double"
		} else {
			return "double"
		}
	case spec.TypeDecimal:
		return "BigDecimal"
	case spec.TypeBoolean:
		if referenceTypesOnly {
			return "Boolean"
		} else {
			return "boolean"
		}
	case spec.TypeString:
		return "String"
	case spec.TypeUuid:
		return "UUID"
	case spec.TypeDate:
		return "LocalDate"
	case spec.TypeDateTime:
		return "LocalDateTime"
	case spec.TypeJson:
		return "Object"
	default:
		return typ
	}
}
