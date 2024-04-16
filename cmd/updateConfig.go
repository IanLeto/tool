package cmd

import (
	"github.com/spf13/cobra"
)

var BatchConfigmap = &cobra.Command{
	Use: "update",
	Run: func(cmd *cobra.Command, args []string) {
		// Do something
	},
}

func init() {
	KubeYaml.Flags().StringP("stdin", "", "", "Read from stdin")
	KubeYaml.Flags().IntP("interval", "", 2, "Interval for something")
	KubeYaml.Flags().IntP("size", "", 10, "Size for something")
}
