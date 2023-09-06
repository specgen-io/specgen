package client

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/goven/spec"
	"strings"
)

func clientName(api *spec.Api) string {
	return fmt.Sprintf(`%sClient`, api.Name.PascalCase())
}

func trimSlash(param string) string {
	return strings.Trim(param, "/")
}
