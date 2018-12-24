package main

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/mattn/go-colorable"

	"github.com/gofunct/hack/pkg/cli"
	"github.com/gofunct/hack/pkg/hackcmd"
	"github.com/gofunct/hack/pkg/hackcmd/cmd"
)

var (
	inReader            = os.Stdin
	outWriter io.Writer = os.Stdout
	errWriter io.Writer = os.Stderr
)

func init() {
	if runtime.GOOS == "windows" {
		outWriter = colorable.NewColorableStdout()
		errWriter = colorable.NewColorableStderr()
	}
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = cmd.NewGrapiCommand(&hackcmd.Ctx{
		IO: &cli.IO{
			In:  inReader,
			Out: outWriter,
			Err: errWriter,
		},
		RootDir: cli.RootDir(cwd),
		Build: hackcmd.BuildConfig{
			AppName:   name,
			Version:   version,
			Revision:  revision,
			BuildDate: buildDate,
			Prebuilt:  prebuilt,
		},
	}).Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
