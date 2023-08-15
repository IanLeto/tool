package cmd

import (
	"context"
	"flag"
	"fmt"
	"github.com/jcmturner/gokrb5/v8/config"
	"github.com/spf13/cobra"
	v12 "k8s.io/api/authentication/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/apis/audit"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

type KubeConn struct {
	ClientSet     kubernetes.Interface
	DynamicClient dynamic.Interface
}

func NewK8sConn(ctx context.Context, conf *config.Config) *KubeConn {
	k8sconfig := flag.String("k8sconfig1", "/Users/ian/.kube/configFromFlags", "kubernetes configFromFlags file path")
	flag.Parse()
	configFromFlags, err := clientcmd.BuildConfigFromFlags("", *k8sconfig)
	if err != nil {
		panic(err)
	}
	client, err := kubernetes.NewForConfig(configFromFlags)
	if err != nil {
		log.Fatal(err)
	}
	dyClient, err := dynamic.NewForConfig(configFromFlags)
	if err != nil {
		log.Fatal(err)
	}
	return &KubeConn{ClientSet: client, DynamicClient: dyClient}
}

var K8sCmd = &cobra.Command{
	Use: "k8s",
	Run: func(cmd *cobra.Command, args []string) {
		conn := NewK8sConn(context.TODO(), nil)
		events, err := conn.ClientSet.CoreV1().Events("default").List(context.TODO(), metav1.ListOptions{})
		NoErr(err)
		for _, item := range events.Items {
			//item.InvolvedObject
			fmt.Println(item.InvolvedObject)
		}
		var event = v1.Event{
			TypeMeta: metav1.TypeMeta{
				Kind:       "",
				APIVersion: "",
			},
			ObjectMeta: metav1.ObjectMeta{ // 事件的标准对象metadata
				Name:                       "",
				GenerateName:               "",
				Namespace:                  "",
				SelfLink:                   "",
				UID:                        "",
				ResourceVersion:            "",
				Generation:                 0,
				CreationTimestamp:          metav1.Time{},
				DeletionTimestamp:          nil,
				DeletionGracePeriodSeconds: nil,
				Labels:                     nil,
				Annotations:                nil,
				OwnerReferences:            nil,
				Finalizers:                 nil,
				ManagedFields:              nil,
			},
			InvolvedObject: v1.ObjectReference{ // 这个事件对应的资源对象
				Kind:            "",
				Namespace:       "",
				Name:            "",
				UID:             "",
				APIVersion:      "",
				ResourceVersion: "",
				FieldPath:       "",
			},
			Reason:  "",
			Message: "",
			Source: v1.EventSource{ //// The component reporting this event. Should be a short machine understandable string.
				Component: "", // 生成事件的组件
				Host:      "", // 生成事件的节点
			},
			FirstTimestamp:      metav1.Time{}, // 事件首次发生时间
			LastTimestamp:       metav1.Time{}, // 事件最后一次发生事件
			Count:               0,
			Type:                "",
			EventTime:           metav1.MicroTime{}, // 该事件首次被观察到的时间
			Series:              nil,
			Action:              "", // 针对该事件所指向对象的相关动作
			Related:             nil,
			ReportingController: "",
			ReportingInstance:   "",
		}
		fmt.Println(event)
		var auditEvent = audit.Event{
			TypeMeta:   metav1.TypeMeta{},
			Level:      "",
			AuditID:    "",
			Stage:      "",
			RequestURI: "",
			Verb:       "",
			User: v12.UserInfo{
				Username: "",
				UID:      "",
				Groups:   nil,
				Extra:    nil,
			},
			ImpersonatedUser: &v12.UserInfo{
				Username: "",
				UID:      "",
				Groups:   nil,
				Extra:    nil,
			},
			SourceIPs:                nil,
			UserAgent:                "",
			ObjectRef:                nil,
			ResponseStatus:           nil,
			RequestObject:            nil,
			ResponseObject:           nil,
			RequestReceivedTimestamp: metav1.MicroTime{}, // 抵达apisever 时间
			StageTimestamp:           metav1.MicroTime{}, //抵达当前审计阶段时间
			Annotations:              nil,
		}
		fmt.Println(auditEvent)
	},
}

func init() {
	K8sCmd.Flags().String("port", "8080", "监听端口")

}
