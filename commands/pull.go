package commands

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/MichaelMure/git-bug/commands/completion"
	"github.com/MichaelMure/git-bug/commands/execenv"
	"github.com/MichaelMure/git-bug/entity"
)

func newPullCommand() *cobra.Command {
	env := execenv.NewEnv()

	cmd := &cobra.Command{
		Use:     "pull [REMOTE]",
		Short:   "Pull updates from a git remote",
		PreRunE: execenv.LoadBackend(env),
		RunE: execenv.CloseBackend(env, func(cmd *cobra.Command, args []string) error {
			return runPull(env, args)
		}),
		ValidArgsFunction: completion.GitRemote(env),
	}

	return cmd
}

func runPull(env *execenv.Env, args []string) error {
	if len(args) > 1 {
		return errors.New("Only pulling from one remote at a time is supported")
	}

	remote := "origin"
	if len(args) == 1 {
		remote = args[0]
	}

	env.Out.Println("Fetching remote ...")

	stdout, err := env.Backend.Fetch(remote)
	if err != nil {
		return err
	}

	env.Out.Println(stdout)

	env.Out.Println("Merging data ...")

	for result := range env.Backend.MergeAll(remote) {
		if result.Err != nil {
			env.Err.Println(result.Err)
		}

		if result.Status != entity.MergeStatusNothing {
			env.Out.Printf("%s: %s\n", result.Id.Human(), result)
		}
	}

	return nil
}
