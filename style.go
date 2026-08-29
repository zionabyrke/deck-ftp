package main

import (
	"os"

	"github.com/fatih/color"
)

var (
	successColor = color.New(color.FgGreen)
	warnColor    = color.New(color.FgYellow)
	errorColor   = color.New(color.FgRed)
)

func success(format string, a ...any) {
	successColor.Printf(format+"\n", a...)
}

func warn(format string, a ...any) {
	warnColor.Printf(format+"\n", a...)
}

// fail prints in red and exits, the styled replacement for log.Fatalf.
func fail(format string, a ...any) {
	errorColor.Printf(format+"\n", a...)
	os.Exit(1)
}
