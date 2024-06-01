package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func NewK8sClient(configPath string) *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		panic(err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return client
}

var AgentCmd = &cobra.Command{
	Use: "agent",
	Run: func(cmd *cobra.Command, args []string) {
		//cmder := exec.Command("/bin/bash", "/home/ian/workdir/tool/cmd/agent_cmd.sh")
		//output, err := cmder.CombinedOutput() // 获取标准输出和标准错误
		//if err != nil {
		//	fmt.Println("Error executing script:", err)
		//	return
		//}
		//var (
		//	client = NewK8sClient("/home/ian/.kube/config")
		//)
		//
		//fmt.Println("Script output:", string(output))
		//switch opt := cmd.Flag("kube").Value.String(); opt {
		//case "kubeconfig":
		//
		//default:
		//
		//}
	},
}

func init() {
	AgentCmd.Flags().StringP("kube", "k", "", "config")
}
