package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"sort"

	"io"
	"os"
)

var DockerCmd = &cobra.Command{
	Use: "docker",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		// 获取命令行参数
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}

		// 监听Docker的事件
		messages, errs := cli.Events(context.Background(), types.EventsOptions{})

		// 使用select监听消息和错误
		containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
		NoErr(err)
		for _, container := range containers {
			fmt.Printf("Container: %v\n", ToJSON(container))
			epl, err := cli.ContainerInspect(ctx, container.ID)

			NoErr(err)
			fmt.Println("容器Info", ToJSON(epl))
			fmt.Println("存储引擎", epl.GraphDriver.Name)
			fmt.Println("存储引擎上级目录", epl.GraphDriver.Data["MergedDir"])
			fmt.Println("存储引擎上级目录", epl.GraphDriver.Data["UpperDir"])
			mounts := map[string]string{}
			for _, mount := range epl.Mounts {
				mounts[mount.Destination] = mount.Source
			}
			mountsKeys := []string{}
			for k := range mounts {
				mountsKeys = append(mountsKeys, k)
			}
			sort.Sort(sort.Reverse(sort.StringSlice(mountsKeys)))
			fmt.Println("挂载点", mountsKeys)

			for _, destination := range mountsKeys {
				//fmt.Printf("挂载点: %v, 挂载源: %v\n", destination, strings.Replace("<pod-id>/project/1.log", mounts[destination], "", 1))
				fmt.Printf("挂载点: %v, 挂载源: %v\n", destination, mounts[destination])

			}

		}

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
					for _, m := range containerJson.Mounts {
						v, _ := json.MarshalIndent(m, "", "  ")
						fmt.Println(string(v))
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
