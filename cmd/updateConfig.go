package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var k8sclient = func() *kubernetes.Clientset {

	config, err := clientcmd.BuildConfigFromFlags("", "/home/ian/.kube/config")
	NoErr(err)
	clientset, err := kubernetes.NewForConfig(config)
	NoErr(err)
	return clientset

}()

// getConfigMap 获取指定的 ConfigMap
func getConfigMap(clientset *kubernetes.Clientset, namespace, name string) (*v1.ConfigMap, error) {
	return clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

// updateConfigMap 更新 ConfigMap
func updateConfigMap(clientset *kubernetes.Clientset, namespace string, configMap *v1.ConfigMap) error {
	_, err := clientset.CoreV1().ConfigMaps(namespace).Update(context.TODO(), configMap, metav1.UpdateOptions{})
	return err
}

// updateQueueSpoolInConfigMap 修改 ConfigMap 中的 YAML 数据
func updateQueueSpoolInConfigMap(cm *v1.ConfigMap, newSize, newBufferSize, newFlushTimeout string) error {
	yamlData, ok := cm.Data["filebeat.yml"]
	if !ok {
		return fmt.Errorf("config.yaml not found in ConfigMap")
	}

	var config map[string]interface{}
	err := yaml.Unmarshal([]byte(yamlData), &config)
	if err != nil {
		return err
	}

	// 修改 queue.spool
	queueSpool, ok := config["queue.spool"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("queue.spool structure not found or incorrect format")
	}

	file, ok := queueSpool["file"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("file structure not found or incorrect format")
	}
	file["size"] = newSize

	write, ok := queueSpool["write"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("write structure not found or incorrect format")
	}
	write["buffer_size"] = newBufferSize
	write["flush.timeout"] = newFlushTimeout

	// 序列化回 YAML
	modifiedData, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	cm.Data["filebeat.yml"] = string(modifiedData)

	return nil
}

var BatchConfigmap = &cobra.Command{
	Use: "update",
	Run: func(cmd *cobra.Command, args []string) {
		// Do something
		l, err := k8sclient.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		NoErr(err)
		for _, ns := range l.Items {
			fmt.Println(ns.Name)
		}
		var (
			configmapname string
			ns            string
		)
		configmapname, _ = cmd.Flags().GetString("configmapname")
		ns, _ = cmd.Flags().GetString("ns")
		cm, err := getConfigMap(k8sclient, ns, configmapname)
		NoErr(err)
		err = updateQueueSpoolInConfigMap(cm, "100M", "1000", "10s")
		NoErr(err)

	},
}

func init() {
	BatchConfigmap.Flags().StringP("config", "", "", "k8s 配置文件路径")
	BatchConfigmap.Flags().StringP("ns", "", "", "k8s 配置文件路径")
	BatchConfigmap.Flags().StringP("configmapname", "c", "", "k8s 配置文件路径")

}
