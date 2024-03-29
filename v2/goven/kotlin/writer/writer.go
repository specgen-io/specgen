package writer

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/goven/generator"
	"github.com/specgen-io/specgen/v2/goven/kotlin/packages"
	"strings"
)

func KotlinConfig() generator.Config {
	return generator.Config{"    ", 4, map[string]string{}}
}

type Writer struct {
	generator.Writer
	thePackage packages.Package
	className  string
	Imports    *imports
}

func New(thePackage packages.Package, className string) *Writer {
	config := KotlinConfig()
	config.Substitutions["[[.ClassName]]"] = className
	return &Writer{generator.NewWriter(config), thePackage, className, NewImports()}
}

func (w *Writer) Indented() *Writer {
	return &Writer{w.Writer.Indented(), w.thePackage, w.className, w.Imports}
}

func (w *Writer) IndentedWith(size int) *Writer {
	return &Writer{w.Writer.IndentedWith(size), w.thePackage, w.className, w.Imports}
}

func (w *Writer) ToCodeFile() *generator.CodeFile {
	lines := []string{
		fmt.Sprintf(`package %s`, w.thePackage.PackageName),
		"",
	}
	imports := w.Imports.Lines()
	if len(imports) > 0 {
		lines = append(lines, imports...)
		lines = append(lines, "")
	}
	lines = append(lines, w.Code()...)
	code := strings.Join(lines, "\n")
	return &generator.CodeFile{w.thePackage.GetPath(fmt.Sprintf("%s.kt", w.className)), code}
}
