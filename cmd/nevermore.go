package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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

func readAndParseYAMLFiles(directory string) {
	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".yaml" || filepath.Ext(file.Name()) == ".yml" {
			filePath := filepath.Join(directory, file.Name())
			data, err := ioutil.ReadFile(filePath)
			if err != nil {
				log.Printf("Failed to read file %s: %v", file.Name(), err)
				continue
			}

			var config FileBeatConfig
			if err := yaml.Unmarshal(data, &config); err != nil {
				log.Printf("Failed to unmarshal YAML file %s: %v", file.Name(), err)
				continue
			}

			fmt.Printf("Parsed YAML file: %s\n", file.Name())
			fmt.Printf("Config: %+v\n", config)
			fmt.Println(config.Output.Kafka.Hosts[0])
		}
	}
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

var NevermoreCmd = &cobra.Command{
	Use: "nevermore",
	Run: func(cmd *cobra.Command, args []string) {
		command, _ := cmd.Flags().GetStringArray("command")
		exportConfigMap, _ := cmd.Flags().GetBool("export-configmap")
		deployConfigMap, _ := cmd.Flags().GetBool("deploy-configmap")
		err := executeCommandsOnClusters(command, exportConfigMap, deployConfigMap)
		cobra.CheckErr(err)
		dir, _ := cmd.Flags().GetString("read-configmap")
		readAndParseYAMLFiles(dir)
	},
}

func init() {

	NevermoreCmd.Flags().StringArrayP("command", "c", []string{}, "Commands to execute on each cluster")
	NevermoreCmd.Flags().BoolP("export-configmap", "e", false, "Export ConfigMap to file")
	NevermoreCmd.Flags().BoolP("deploy-configmap", "d", false, "Deploy ConfigMap from file")
	NevermoreCmd.Flags().BoolP("read-configmap", "", false, "read当前目录的yaml文件")
	// nevermore  --export-configmap
	// nevermore  --read-configmap /tmp
}
