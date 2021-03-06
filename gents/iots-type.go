package gents

import (
	"fmt"
	"github.com/specgen-io/spec"
)

func IoTsType(typ *spec.TypeDef) string {
	switch typ.Node {
	case spec.PlainType:
		return PlainIoTsType(typ.Plain)
	case spec.NullableType:
		child := IoTsType(typ.Child)
		result := "t.union([" + child + ", t.null])"
		return result
	case spec.ArrayType:
		child := IoTsType(typ.Child)
		result := "t.array(" + child + ")"
		return result
	case spec.MapType:
		child := IoTsType(typ.Child)
		result := "t.record(t.string, " + child + ")"
		return result
	default:
		panic(fmt.Sprintf("Unknown type: %v", typ))
	}
}

func PlainIoTsType(typ string) string {
	switch typ {
	case spec.TypeInt32:
		return "t.number"
	case spec.TypeInt64:
		return "t.number"
	case spec.TypeFloat:
		return "t.number"
	case spec.TypeDouble:
		return "t.number"
	case spec.TypeDecimal:
		return "t.number"
	case spec.TypeBoolean:
		return "t.boolean"
	case spec.TypeString:
		return "t.string"
	case spec.TypeUuid:
		return "t.string"
	case spec.TypeDate:
		return "t.string"
	case spec.TypeDateTime:
		return "t.DateFromISOString"
	case spec.TypeJson:
		return "t.unknown"
	default:
		return "T" + typ
	}
}
