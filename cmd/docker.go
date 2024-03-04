package cmd

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"io"
	"os"
)

var DockerCmd = &cobra.Command{
	Use: "docker",
	Run: func(cmd *cobra.Command, args []string) {
		// 获取命令行参数
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}

		// 监听Docker的事件
		messages, errs := cli.Events(context.Background(), types.EventsOptions{})

		// 使用select监听消息和错误
		cli.ContainerList(context.Background(), types.ContainerListOptions{})
		for {
			select {
			case err := <-errs:
				if err != nil && err != io.EOF {
					fmt.Fprintf(os.Stderr, "error: %v\n", err)
					os.Exit(1)
				}
				return
			case msg := <-messages:
				if msg.Type == "container" {
					containerJson, err := cli.ContainerInspect(context.Background(), msg.Actor.ID)
					if err != nil {
						fmt.Sprintf("Error: %v\n", err)
					}
					fmt.Printf("Container: %v\n", containerJson)

				}

			}
		}
	},
}

func init() {
	DockerCmd.Flags().StringP("address", "", "", "")

}
