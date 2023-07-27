package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

var CronJobCmd = &cobra.Command{
	Use: "cron",
	Run: func(cmd *cobra.Command, args []string) {
		interval, _ := cmd.Flags().GetInt("interval")
		size, _ := cmd.Flags().GetInt("size")
		ticket := time.NewTicker(time.Duration(interval) * time.Second)
		var i = 1
		for {
			select {
			case <-ticket.C:
				for i := 0; i < size; i++ {
					fmt.Println(i)
				}
				i += 1
			}
		}

	},
}

func init() {
	CronJobCmd.Flags().StringP("stdin", "", "", "")
	CronJobCmd.Flags().IntP("interval", "", 2, "")
	CronJobCmd.Flags().IntP("size", "", 10, "")
}
