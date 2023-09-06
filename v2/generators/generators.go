package generators

import (
	"github.com/specgen-io/specgen/v2/goven/generator"
	golang "github.com/specgen-io/specgen/v2/goven/golang/generators"
	java "github.com/specgen-io/specgen/v2/goven/java/generators"
	kotlin "github.com/specgen-io/specgen/v2/goven/kotlin/generators"
	"github.com/specgen-io/specgen/v2/goven/openapi"
	ruby "github.com/specgen-io/specgen/v2/goven/ruby/generators"
	scala "github.com/specgen-io/specgen/v2/goven/scala/generators"
	typescript "github.com/specgen-io/specgen/v2/goven/typescript/generators"
)

var All = []generator.Generator{
	golang.Models,
	golang.Client,
	golang.Service,
	java.Models,
	java.Client,
	java.Service,
	kotlin.Models,
	kotlin.Client,
	kotlin.Service,
	ruby.Models,
	ruby.Client,
	scala.Models,
	scala.Client,
	scala.Service,
	typescript.Models,
	typescript.Client,
	typescript.Service,
	openapi.Openapi,
}
