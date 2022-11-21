package cmd

import (
	"github.com/specgen-io/specgen/v2/generators"
	"github.com/specgen-io/specgen/v2/goven/generator"
	"github.com/spf13/cobra"
)

func init() {
	generator.AddCobraCommands(cmdCodegen, generators.All)
	rootCmd.AddCommand(cmdCodegen)
}

var cmdCodegen = &cobra.Command{
	Use:   "codegen",
	Short: "Generate code from spec using specgen code generators",
}
