package beardsrv

import (
	"fmt"

	"github.com/spf13/cobra"
)

const cmdName = "beardsrv"

var cmdRoot = &cobra.Command{
	Use:     cmdName,
	Version: "v0.0.1",
	RunE:    runRoot,
}

func runRoot(cmd *cobra.Command, args []string) error {
	fmt.Println("butts.")
	return nil
}

func Run() error {
	return cmdRoot.Execute()
}
