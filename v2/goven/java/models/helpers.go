package models

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/goven/spec"
	"strings"
)

func getterName(field *spec.NamedDefinition) string {
	return fmt.Sprintf(`get%s`, field.Name.PascalCase())
}

func setterName(field *spec.NamedDefinition) string {
	return fmt.Sprintf(`set%s`, field.Name.PascalCase())
}

func joinParams(params []string) string {
	return strings.Join(params, ", ")
}
