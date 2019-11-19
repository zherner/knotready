package main

import (
	"flag"
	"fmt"
	"strings"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	 _ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	var kubeconfig *string

    if kubeconfig_envvar := os.Getenv("KUBECONFIG"); kubeconfig_envvar != "" {
        kubeconfig_tokenized := strings.Split(kubeconfig_envvar, ";")
        kubeconfig = flag.String("kubeconfig", kubeconfig_tokenized[0], "absolute path to the kubeconfig file")
    } else if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods\n", len(pods.Items))
	nonReadyPods := 0
	fmt.Printf("These pods are not ready: \n")
	for _, pod := range pods.Items {
		if(pod.Status.Phase != v1.PodRunning && pod.Status.Phase != v1.PodSucceeded) {
			fmt.Printf("Name: %s - Status: %s - Reason: %s \n", pod.Name, pod.Status.Phase, pod.Status.Reason)
			nonReadyPods++
		}
	}
	fmt.Printf("There are %d non-ready pods\n", nonReadyPods)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

