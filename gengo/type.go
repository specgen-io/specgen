package gengo

import (
	"fmt"
	spec "github.com/specgen-io/spec"
)

func GoType(typ *spec.TypeDef) string {
	switch typ.Node {
	case spec.PlainType:
		return PlainGoType(typ.Plain)
	case spec.NullableType:
		child := GoType(typ.Child)
		result := "*" + child
		return result
	case spec.ArrayType:
		child := GoType(typ.Child)
		result := "[]" + child
		return result
	case spec.MapType:
		child := GoType(typ.Child)
		result := "map[string]" + child
		return result
	default:
		panic(fmt.Sprintf("Unknown type: %v", typ))
	}
}

func PlainGoType(typ string) string {
	switch typ {
	case spec.TypeInt32:
		return "int"
	case spec.TypeInt64:
		return "int64"
	case spec.TypeFloat:
		return "float32"
	case spec.TypeDouble:
		return "float64"
	case spec.TypeDecimal:
		return "decimal.Decimal"
	case spec.TypeBoolean:
		return "bool"
	case spec.TypeString:
		return "string"
	case spec.TypeUuid:
		return "uuid.UUID"
	case spec.TypeDate:
		return "civil.Date"
	case spec.TypeDateTime:
		return "civil.DateTime"
	case spec.TypeJson:
		return "json.RawMessage"
	case spec.TypeEmpty:
		return "EmptyDef"
	default:
		return typ
	}
}
