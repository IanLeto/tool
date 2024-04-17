package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

var client = func() *kubernetes.Clientset {
	kubeconfig := os.Getenv("KUBECONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	NoErr(err)
	clientset, err := kubernetes.NewForConfig(config)
	NoErr(err)
	return clientset

}()

var BatchConfigmap = &cobra.Command{
	Use: "update",
	Run: func(cmd *cobra.Command, args []string) {
		// Do something

	},
}

func init() {
	KubeYaml.Flags().StringP("stdin", "", "", "Read from stdin")
}
