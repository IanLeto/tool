package cmd

import (
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
)

type KubeConfig struct {
	APIVersion string `yaml:"apiVersion"`
	Clusters   []struct {
		Cluster struct {
			InsecureSkipTLSVerify bool   `yaml:"insecure-skip-tls-verify"`
			Server                string `yaml:"server"`
		} `yaml:"cluster"`
		Name string `yaml:"name"`
	} `yaml:"clusters"`
	Contexts []struct {
		Context struct {
			Cluster   string `yaml:"cluster"`
			User      string `yaml:"user"`
			Namespace string `yaml:"namespace"`
		} `yaml:"context"`
		Name string `yaml:"name"`
	} `yaml:"contexts"`
	CurrentContext string   `yaml:"current-context"`
	Kind           string   `yaml:"kind"`
	Preferences    struct{} `yaml:"preferences"`
	Users          []struct {
		Name string `yaml:"name"`
		User struct {
			Token string `yaml:"token"`
		} `yaml:"user"`
	} `yaml:"users"`
}

func mergeKubeConfigs(files []string) (*KubeConfig, error) {
	var mergedConfig KubeConfig
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}

		var config KubeConfig
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, err
		}

		// Merge logic
		mergedConfig.Clusters = append(mergedConfig.Clusters, config.Clusters...)
		mergedConfig.Contexts = append(mergedConfig.Contexts, config.Contexts...)
		mergedConfig.Users = config.Users // Assuming all users are the same
	}

	// Assuming the first file sets the general structure
	//if len(files) > 0 {
	//	firstConfigData, _ := ioutil.ReadFile(files[0])
	//	err := yaml.Unmarshal(firstConfigData, &mergedConfig)
	//	NoErr(err)
	//}

	return &mergedConfig, nil
}

func traverseAndMerge(rootDir, outputFile string) error {
	var files []string
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	mergedConfig, err := mergeKubeConfigs(files)
	if err != nil {
		return err
	}

	outData, err := yaml.Marshal(mergedConfig)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(outputFile, outData, 0644); err != nil {
		return err
	}

	return nil
}

var KubeYaml = &cobra.Command{
	Use: "KubeYaml",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")
		outputFile, _ := cmd.Flags().GetString("output")
		err := traverseAndMerge(dir, outputFile)
		NoErr(err)
	},
}

func init() {
	KubeYaml.Flags().StringP("dir", "", "", "kube 目录")
	KubeYaml.Flags().StringP("output", "", "", "输出文件")

}

// It's important to note that you need to replace the initialization of the rest.Config
// with actual values and also adjust the convertToKubeConfig function according to your
// authentication and cluster setup. This example provides a starting point for you to
// customize the conversion to your specific needs.
