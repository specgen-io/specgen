package cmd

import (
	"github.com/specgen-io/specgen/v2/goven/generator/console"
	"os"
)

func Fail(args ...interface{}) {
	console.ProblemLn(args...)
	os.Exit(1)
}

func FailF(format string, args ...interface{}) {
	console.ProblemLnF(format, args...)
	os.Exit(1)
}

func FailIfError(err error, args ...interface{}) {
	if err != nil {
		console.ProblemLn(args...)
		console.ProblemLn(err)
		os.Exit(1)
	}
}

func FailIfErrorF(err error, format string, args ...interface{}) {
	if err != nil {
		console.ProblemLnF(format, args...)
		console.ProblemLn(err)
		os.Exit(1)
	}
}
