package empty

import (
	"github.com/specgen-io/specgen/v2/goven/generator"
	"github.com/specgen-io/specgen/v2/goven/golang/module"
	"github.com/specgen-io/specgen/v2/goven/golang/writer"
)

func GenerateEmpty(emptyModule module.Module) *generator.CodeFile {
	w := writer.New(emptyModule, `empty.go`)
	w.Lines(`
type Type struct{}

var Value = Type{}
`)
	return w.ToCodeFile()
}
