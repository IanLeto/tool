package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"k8s.io/client-go/tools/clientcmd/api"
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
	InsecureSkipTLSVerify    bool   `yaml:"-"`
	CertificateAuthorityData []byte `yaml:"-"`
	// 如果有其他字段也想忽略，也可以加上 `yaml:"-"`
}

type CustomContext struct {
	Cluster   string `yaml:"cluster"`
	AuthInfo  string `yaml:"authinfo"`
	Namespace string `yaml:"namespace"`
	// 忽略不需要的字段
	// Extensions               map[string]interface{} `yaml:"-"`
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
			Clusters       map[string]CustomCluster `yaml:"clusters"`
			AuthInfos      map[string]*api.AuthInfo `yaml:"authinfos"`
			Contexts       map[string]CustomContext `yaml:"contexts"`
			CurrentContext string                   `yaml:"currentcontext"`
		}{
			Kind:           "config",
			APIVersion:     "v1",
			Clusters:       map[string]CustomCluster{},
			AuthInfos:      nil,
			Contexts:       map[string]CustomContext{},
			CurrentContext: "",
		}

		for _, item := range clusters.Items {
			clusterName := item.Metadata.Name
			kubeconfig.Clusters[clusterName] = CustomCluster{
				Server:                   "https://api." + clusterName + ".com",
				CertificateAuthorityData: []byte(""),
			}
			kubeconfig.Contexts[clusterName] = CustomContext{
				Cluster:   clusterName,
				AuthInfo:  "k0110",
				Namespace: "cpaas-system",
			}
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
