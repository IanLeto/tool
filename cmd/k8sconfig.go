package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type Clusters struct {
	Items []struct {
		Metadata struct {
			Name string `json:"name"`
		} `json:"metadata"`
	} `json:"items"`
}

type CustomCluster struct {
	Server string `yaml:"server"`
	// 其他字段你不希望显示的可以使用 `yaml:"-"` 来忽略
	InsecureSkipTLSVerify    bool   `yaml:"insecure-skip-tls-verify"`
	CertificateAuthorityData []byte `yaml:"-"`
	// 如果有其他字段也想忽略，也可以加上 `yaml:"-"`
}

type CustomContext struct {
	Cluster   string `yaml:"cluster"`
	User      string `yaml:"user"`
	Namespace string `yaml:"namespace"`
	// 忽略不需要的字段
	// Extensions               map[string]interface{} `yaml:"-"`
}
type SpecCluster struct {
	Name    string        `yaml:"name"`
	Cluster CustomCluster `yaml:"cluster"`
}

type SpecContext struct {
	Name    string        `yaml:"name"`
	Context CustomContext `yaml:"context"`
}

var KubeYaml = &cobra.Command{
	Use: "KubeYaml",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			clusters Clusters
		)
		cluster, err := ioutil.ReadFile("clusters.json")
		if err != nil {
			fmt.Printf("Error reading clusters.json: %v\n", err)
			return
		}

		err = json.Unmarshal(cluster, &clusters)
		if err != nil {
			fmt.Printf("Error unmarshaling clusters: %v\n", err)
			return
		}

		// 使用自定义的结构体
		kubeconfig := struct {
			Kind           string                   `yaml:"kind"`
			APIVersion     string                   `yaml:"apiversion"`
			Clusters       []SpecCluster            `yaml:"clusters"`
			AuthInfos      []map[string]interface{} `yaml:"users"`
			Contexts       []SpecContext            `yaml:"contexts"`
			CurrentContext string                   `yaml:"current-context"`
		}{
			APIVersion:     "v1",
			Kind:           "config",
			Clusters:       []SpecCluster{},
			AuthInfos:      []map[string]interface{}{},
			Contexts:       []SpecContext{},
			CurrentContext: "",
		}

		for _, item := range clusters.Items {
			clusterName := item.Metadata.Name
			kubeconfig.Clusters = append(kubeconfig.Clusters, SpecCluster{
				Name: clusterName,
				Cluster: CustomCluster{
					Server: "https://api." + clusterName + ".com",
				},
			})
			kubeconfig.Contexts = append(kubeconfig.Contexts, SpecContext{
				Name: clusterName,
				Context: CustomContext{
					Cluster:   clusterName,
					User:      "xx",
					Namespace: "default",
				},
			})
		}

		// 省略的 AuthInfo 结构体初始化...

		yamlData, err := yaml.Marshal(&kubeconfig)
		if err != nil {
			fmt.Printf("Error marshaling kubeconfig: %v\n", err)
			return
		}

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
