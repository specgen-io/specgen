package writer

import (
	"github.com/specgen-io/specgen/v2/goven/generator"
)

var RubyConfig = generator.Config{"  ", 2, nil}

func New(filename string) generator.Writer {
	return generator.NewWriter(filename, RubyConfig)
}
