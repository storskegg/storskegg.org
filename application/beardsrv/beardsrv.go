package beardsrv

import (
	"github.com/spf13/cobra"
	"github.com/storskegg/storskegg.org/internal/server"
)

const cmdName = "beardsrv"

var cmdRoot = &cobra.Command{
	Use:     cmdName,
	Version: "v0.0.1",
	RunE:    runRoot,
}

func runRoot(cmd *cobra.Command, args []string) error {
	srv := server.New(&server.Config{Addr: ":3001"}, cmd)
	return srv.Serve()
}

func Run() error {
	return cmdRoot.Execute()
}
