package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
)

type FileBeatConfig struct {
	Filebeat struct {
		Inputs []struct {
			Type            string            `yaml:"type"`
			Enabled         bool              `yaml:"enabled"`
			Paths           []string          `yaml:"paths"`
			Fields          map[string]string `yaml:"fields"`
			FieldsUnderRoot bool              `yaml:"fields_under_root"`
			Encoding        string            `yaml:"encoding"`
			ExcludeLines    []string          `yaml:"exclude_lines"`
		} `yaml:"inputs"`
		Config struct {
			Modules struct {
				Path   string `yaml:"path"`
				Reload struct {
					Enabled bool   `yaml:"enabled"`
					Period  string `yaml:"period"`
				} `yaml:"reload"`
			} `yaml:"modules"`
		} `yaml:"config"`
	} `yaml:"filebeat"`
	Output struct {
		Kafka struct {
			Hosts []string `yaml:"hosts"`
			Topic string   `yaml:"topic"`
		} `yaml:"kafka"`
	} `yaml:"output"`
}

func getClusterNames() ([]string, error) {
	cmd := exec.Command("kubectl", "config", "get-contexts", "-o", "name")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster names: %v", err)
	}

	var clusterNames []string
	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		clusterNames = append(clusterNames, scanner.Text())
	}

	return clusterNames, nil
}

func executeCommandsOnClusters(commands []string) error {
	clusters, err := getClusterNames()
	if err != nil {
		return err
	}

	for _, cluster := range clusters {
		fmt.Printf("Switching to cluster: %s\n", cluster)

		// 切换到指定的集群
		cmd := exec.Command("kubectl", "config", "use-context", cluster)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to switch to cluster %s: %v", cluster, err)
		}
		// 执行命令列表
		for _, command := range commands {
			fmt.Printf("Executing command: %s\n", command)

			// 按逗号分割命令
			cmdList := strings.Split(command, ",")

			for _, cmd := range cmdList {
				// 分割命令和参数

				// 执行命令
				execCmd := exec.Command(cmd)
				execCmd.Stdout = os.Stdout
				execCmd.Stderr = os.Stderr
				err := execCmd.Run()
				if err != nil {
					return fmt.Errorf("failed to execute command %s on cluster %s: %v", cmd, cluster, err)
				}
			}
		}

		fmt.Println("----")
	}

	return nil
}

var ClusterCmd = &cobra.Command{
	Use: "clusters",
	Run: func(cmd *cobra.Command, args []string) {
		command, _ := cmd.Flags().GetStringArray("command")
		err := executeCommandsOnClusters(command)
		cobra.CheckErr(err)
	},
}

func init() {
	ClusterCmd.Flags().StringArrayP("command", "c", []string{}, "Commands to execute on each cluster")
}
