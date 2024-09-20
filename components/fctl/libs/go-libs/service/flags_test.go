package service

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/cobra"
)

func TestFlags(t *testing.T) {

	os.Setenv("ROOT1", "changed")

	command := &cobra.Command{
		Use: "root",
	}
	command.Flags().String("root1", "test", "")

	subCommand1 := &cobra.Command{
		Use: "subcommand1",
	}
	subCommand1.Flags().String("sub1", "test", "")

	subCommand2 := &cobra.Command{
		Use: "subcommand2",
	}
	subCommand2.Flags().String("sub2", "test", "")
	subCommand2.PersistentFlags().String("persub2", "test", "")

	command.AddCommand(subCommand1, subCommand2)

	BindEnvToCommand(command)

	command.Usage()
	fmt.Println(command.Flags().GetString("root1"))
}
