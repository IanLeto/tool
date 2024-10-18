package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
)

func runShellScript(scriptPath string) (string, error) {
	// 使用bash -c 执行脚本
	cmd := exec.Command("bash", scriptPath)

	// 获取命令的输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute script: %v, output: %s", err, string(output))
	}

	// 返回输出结果，去除可能的额外空白字符
	return strings.TrimSpace(string(output)), nil
}

func runKubectlCommand(args []string, kubeconfig string) (string, error) {
	// 设置KUBECONFIG环境变量
	os.Setenv("KUBECONFIG", kubeconfig)

	// 创建kubectl命令
	cmd := exec.Command("kubectl", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute kubectl command: %v, output: %s", err, string(output))
	}
	return string(output), nil
}

var ExecCmd = &cobra.Command{
	Use: "exec",
	Run: func(cmd *cobra.Command, args []string) {
		//// 获取Kubeconfig文件路径
		//kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
		//if configPath := os.Getenv("KUBECONFIG"); configPath != "" {
		//	kubeconfig = configPath
		//}
		//
		//// 加载Kubeconfig
		//config, err := clientcmd.LoadFromFile(kubeconfig)
		//if err != nil {
		//	fmt.Println("Error loading kubeconfig:", err)
		//	return
		//}
	},
}

func init() {
	ExecCmd.Flags().StringP("input", "i", "", "k8s 配置文件")

}
