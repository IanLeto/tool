package cmd

import (
	"context"
	"flag"
	"github.com/spf13/cobra"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	typ "k8s.io/apimachinery/pkg/types"
	addScheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	setupLog = ctrl.Log.WithName("setup")
)

type KubeConn struct {
	ClientSet     kubernetes.Interface
	DynamicClient dynamic.Interface
}

// NewK8sConn initializes a connection to a Kubernetes cluster
func NewK8sConn(ctx context.Context) *KubeConn {
	k8sconfig := flag.String("k8sconfig", "/Users/ian/.kube/config", "Path to the Kubernetes config file")
	flag.Parse()

	configFromFlags, err := clientcmd.BuildConfigFromFlags("", *k8sconfig)
	if err != nil {
		log.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	client, err := kubernetes.NewForConfig(configFromFlags)
	if err != nil {
		log.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	dyClient, err := dynamic.NewForConfig(configFromFlags)
	if err != nil {
		log.Fatalf("Error building dynamic clientset: %s", err.Error())
	}

	return &KubeConn{ClientSet: client, DynamicClient: dyClient}
}

var (
	scheme = runtime.NewScheme()

	// K8sCmd represents the base command when called without any subcommands
	K8sCmd = &cobra.Command{
		Use:   "k8s",
		Short: "K8s is a CLI for managing Kubernetes clusters",
		Long: `K8s is a CLI application for managing Kubernetes clusters,
providing user-friendly interactions for complex operations.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
			log.Println("K8s command called")
			ctx := context.Background()
			ctrl.Log.WithName("controller-runtime").Info("K8s command called")

			// Setup a Manager
			mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
				Scheme: runtime.NewScheme(),
			})
			NoErr(err)
			client := mgr.GetClient()
			err = client.Get(ctx, typ.NamespacedName{Name: os.Getenv("NODE_NAME")}, &corev1.Node{})
			NoErr(err)

			if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
				setupLog.Error(err, "problem running manager")
				return
			}

		},
	}
)

func init() {
	// Initialize scheme
	if err := addScheme.AddToScheme(scheme); err != nil {
		log.Fatalf("Error adding to scheme: %s", err.Error())
	}

	// Define flags for the K8sCmd
	K8sCmd.PersistentFlags().String("k8sconfig", "/Users/ian/.kube/config", "Path to the Kubernetes config file")
	K8sCmd.Flags().StringP("port", "p", "8080", "Port to listen on")
	K8sCmd.Flags().StringP("address", "a", "", "Address to bind to")

	// Add the K8sCmd to the Cobra root command
	// Here you would add the root command if you have one
	// For example, if you have a rootCmd representing the entry point of your application, you would add K8sCmd to it:
	// rootCmd.AddCommand(K8sCmd)
}
