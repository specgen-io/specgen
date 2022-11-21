package writer

import (
	"fmt"
	"github.com/specgen-io/specgen/v2/goven/generator"
	"github.com/specgen-io/specgen/v2/goven/java/packages"
)

func JavaConfig() generator.Config {
	return generator.Config{"\t", 2, map[string]string{}}
}

func New(thePackage packages.Package, className string) generator.Writer {
	config := JavaConfig()
	filename := thePackage.GetPath(fmt.Sprintf("%s.java", className))
	config.Substitutions["[[.ClassName]]"] = className
	w := generator.NewWriter(filename, config)
	w.Line(`package %s;`, thePackage.PackageName)
	w.EmptyLine()
	return w
}
