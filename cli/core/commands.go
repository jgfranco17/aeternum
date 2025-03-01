package core

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	failFast bool
)

func GetPingCommand() *cobra.Command {
	var timeout int
	var count int
	cmd := &cobra.Command{
		Use:   "ping",
		Short: "Ping a target URL",
		Long:  "Check if a target URL is live and responds with a 2xx status code",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("Not enough arguments, expected 1 but got %d", len(args))
			}
			target := args[0]
			if err := PingUrl(target, count, timeout); err != nil {
				log.Errorf("Ping failed: %v", err)
				return err
			}
			return nil
		},
	}
	cmd.Flags().IntVarP(&count, "count", "c", 1, "Number of ping requests, default is 1")
	cmd.Flags().IntVarP(&timeout, "timeout", "t", 5, "Timeout duration (in seconds) for the ping request, default is 5s")
	return cmd
}
