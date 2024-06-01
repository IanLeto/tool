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
	},
}

func init() {
	AgentCmd.Flags().StringP("kube", "k", "", "config")
}
