package models

import (
	"github.com/specgen-io/specgen/v2/goven/spec"
)

func oneOfItemClassName(item *spec.NamedDefinition) string {
	return item.Name.PascalCase()
}
