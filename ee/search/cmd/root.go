package cmd

import (
	_ "github.com/bombsimon/logrusr/v3"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/spf13/cobra"
)

var (
	Version   = "develop"
	BuildDate = "-"
	Commit    = "-"
)

func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:               "search",
		Short:             "search",
		DisableAutoGenTag: true,
		Version:           Version,
	}

	serverCmd := NewServer()
	root.AddCommand(NewVersion(), serverCmd, NewInitMapping(), NewUpdateMapping())

	return root
}

func Execute() {
	service.Execute(NewRootCommand())
}
