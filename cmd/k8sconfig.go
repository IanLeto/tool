package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"k8s.io/client-go/tools/clientcmd/api"
	"os"
)

type Cluster struct {
	Metadata struct {
		Name string `json:"name"`
	} `json:"metadata"`
}

var KubeYaml = &cobra.Command{
	Use: "KubeYaml",
	Run: func(cmd *cobra.Command, args []string) {
		kubeconfig := api.Config{
			Kind:           "config",
			APIVersion:     "v1",
			Clusters:       nil,
			AuthInfos:      nil,
			Contexts:       nil,
			CurrentContext: "",
			Extensions:     nil,
		}
		input, err := os.ReadFile("cluster.json")
		NoErr(err)
		// Marshal the kubeconfig object into YAML
		yamlData, err := yaml.Marshal(&kubeconfig)
		if err != nil {
			fmt.Printf("Error marshaling kubeconfig: %v\n", err)
			return
		}

		// Print to console
		fmt.Println(string(yamlData))

		// Output to a file
		fileName := "kubeconfig.yaml"
		err = os.WriteFile(fileName, yamlData, 0644)
		if err != nil {
			fmt.Printf("Error writing kubeconfig to file: %v\n", err)
			return
		}

		fmt.Printf("kubeconfig is written to %s\n", fileName)
	},
}

func init() {
	KubeYaml.Flags().StringP("stdin", "", "", "Read from stdin")
	KubeYaml.Flags().IntP("interval", "", 2, "Interval for something")
	KubeYaml.Flags().IntP("size", "", 10, "Size for something")
}

// It's important to note that you need to replace the initialization of the rest.Config
// with actual values and also adjust the convertToKubeConfig function according to your
// authentication and cluster setup. This example provides a starting point for you to
// customize the conversion to your specific needs.
