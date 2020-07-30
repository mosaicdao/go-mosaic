package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mosaicdao/go-mosaic/node"
)

func AddNodeFlags(cmd *cobra.Command) {
	cmd.Flags().String("threads.host", config.Threads.HostAddressString, "Host address")
}


func NewRunNodeCmd(nodeProvider node.NodeProvider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node",
		Short: "Run a mosaic node",
		RunE: func(cmd *cobra.Command, args []string) error {
			// check config
			n, err := nodeProvider(config)
			if err != nil {
				return fmt.Errorf("failed to create new node: %w", err)
			}

			if err := n.Start(); err != nil {
				return fmt.Errorf("failed to start node: %w", err)
			}

			log.Info("Started node")

			// TODO: catch SIGTERM interrupt

			n.Wait()
			return nil
		},
	}

	AddNodeFlags(cmd)
	return cmd
}
