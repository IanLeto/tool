package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

var CronJobCmd = &cobra.Command{
	Use: "cron",
	Run: func(cmd *cobra.Command, args []string) {
		ticket := time.NewTicker(2 * time.Second)
		var i = 1
		for {
			select {
			case <-ticket.C:
				fmt.Println(i)
				i += 1
			}
		}

	},
}

func init() {
	CronJobCmd.Flags().StringP("stdin", "", "", "")
}
