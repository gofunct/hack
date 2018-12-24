package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/gofunct/hack/pkg/cli"
	"github.com/gofunct/hack/pkg/hackcmd"
	"github.com/gofunct/hack/pkg/hackcmd/di"
	"github.com/gofunct/hack/pkg/hackcmd/internal/module"
)

func newUserDefinedCommands(ctx *hackcmd.Ctx) (cmds []*cobra.Command) {
	if !ctx.IsInsideApp() {
		return
	}

	scriptLoader := di.NewScriptLoader(ctx)

	err := scriptLoader.Load(ctx.RootDir.Join("cmd"))
	if err != nil {
		// TODO: log
		return
	}

	ui := di.NewUI(ctx)

	for _, name := range scriptLoader.Names() {
		cmds = append(cmds, newUserDefinedCommand(ui, scriptLoader, name))
	}

	return
}

func newUserDefinedCommand(ui cli.UI, scriptLoader module.ScriptLoader, name string) *cobra.Command {
	return &cobra.Command{
		Use:           name + " [-- BUILD_OPTIONS] [-- RUN_ARGS]",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(c *cobra.Command, args []string) (err error) {
			script, ok := scriptLoader.Get(name)
			if !ok {
				err = errors.Wrapf(err, "faild to find subcommand %s", name)
				return
			}

			pos := len(args)
			for i, arg := range args {
				if arg == "--" {
					pos = i
					break
				}
			}
			var buildArgs, runArgs []string
			if pos == len(args) {
				buildArgs = args
			} else {
				buildArgs = args[:pos]
				runArgs = args[pos+1:]
			}

			ui.Section(script.Name())
			ui.Subsection("Building...")
			err = errors.WithStack(script.Build(buildArgs...))
			if err != nil {
				return
			}

			ui.Subsection("Starting...")
			return errors.WithStack(script.Run(runArgs...))
		},
	}
}
