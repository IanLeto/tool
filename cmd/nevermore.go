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

func executeCommandsOnClusters(commands []string, exportConfigMap bool, deployConfigMap bool) error {
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

		if exportConfigMap {
			// 导出 ConfigMap 为文件
			configMapName := "log-agent-template"
			fileName := fmt.Sprintf("%s-%s.yaml", cluster, configMapName)
			cmd := exec.Command("kubectl", "get", "configmap", configMapName, "-o", "yaml")
			output, err := cmd.Output()
			if err != nil {
				fmt.Printf("failed to export ConfigMap %s from cluster %s: %v", configMapName, cluster, err)
				break
			}
			err = os.WriteFile(fileName, output, 0644)
			if err != nil {
				return fmt.Errorf("failed to write ConfigMap file %s: %v", fileName, err)
			}
			fmt.Printf("Exported ConfigMap %s from cluster %s to file %s\n", configMapName, cluster, fileName)
		}

		if deployConfigMap {
			// 部署 ConfigMap 到集群
			fileName := fmt.Sprintf("%s-log-agent-template.yaml", cluster)
			cmd := exec.Command("kubectl", "apply", "-f", fileName)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				return fmt.Errorf("failed to deploy ConfigMap file %s to cluster %s: %v", fileName, cluster, err)
			}
			fmt.Printf("Deployed ConfigMap file %s to cluster %s\n", fileName, cluster)
		}

		// 执行命令列表
		for _, command := range commands {
			fmt.Printf("Executing command: %s\n", command)

			// 按逗号分割命令
			cmdList := strings.Split(command, ",")

			for _, cmd := range cmdList {
				// 分割命令和参数
				cmdParts := strings.Fields(cmd)
				if len(cmdParts) == 0 {
					continue
				}

				// 执行命令
				execCmd := exec.Command(cmdParts[0], cmdParts[1:]...)
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

var NevermoreCmd = &cobra.Command{
	Use: "nevermore",
	Run: func(cmd *cobra.Command, args []string) {
		command, _ := cmd.Flags().GetStringArray("command")
		exportConfigMap, _ := cmd.Flags().GetBool("export-configmap")
		deployConfigMap, _ := cmd.Flags().GetBool("deploy-configmap")
		err := executeCommandsOnClusters(command, exportConfigMap, deployConfigMap)
		cobra.CheckErr(err)
	},
}

func init() {
	NevermoreCmd.Flags().StringArrayP("command", "c", []string{}, "Commands to execute on each cluster")
	NevermoreCmd.Flags().BoolP("export-configmap", "e", false, "Export ConfigMap to file")
	NevermoreCmd.Flags().BoolP("deploy-configmap", "d", false, "Deploy ConfigMap from file")
}
