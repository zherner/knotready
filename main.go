/*
Copyright 2016 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Note: the example only works with the code within the same release/branch.
package main

import (
    "flag"
    "fmt"
    "strconv"
    "strings"
    "os"
    "path/filepath"

    "k8s.io/apimachinery/pkg/api/errors"
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
    fmt.Printf("Number of pods registered: %d\n", len(pods.Items))
    nonReadyPods := 0
    for index, pod := range pods.Items {
        if(pod.Status.Phase != v1.PodRunning) {
            fmt.Printf("tmp")
            fmt.Printf(strconv.Itoa(index))
            nonReadyPods++
        }
    }
    fmt.Printf("Number of non-running pods: %d\n", nonReadyPods)

    deploys, err := clientset.AppsV1().Deployments("").List(metav1.ListOptions{})
    if err != nil {
        panic(err.Error())
    }
    fmt.Printf("Number of deployments registered: %d\n", len(deploys.Items))
    nonReadyDeploys := 0
    for index, deploy := range deploys.Items {
        if(deploy.Status.UnavailableReplicas > 0) {
            fmt.Printf("tmp")
            fmt.Printf(strconv.Itoa(index))
            nonReadyDeploys++
        }
    }
    fmt.Printf("Nummber of incomplete deployments: %d\n", nonReadyDeploys)

    // Examples for error handling:
    // - Use helper functions like e.g. errors.IsNotFound()
    // - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
    namespace := "default"
    pod := "example-xxxxx"
    _, err = clientset.CoreV1().Pods(namespace).Get(pod, metav1.GetOptions{})
    if errors.IsNotFound(err) {
        fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
    } else if statusError, isStatus := err.(*errors.StatusError); isStatus {
        fmt.Printf("Error getting pod %s in namespace %s: %v\n",
            pod, namespace, statusError.ErrStatus.Message)
    } else if err != nil {
        panic(err.Error())
    } else {
        fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
    }

}

func homeDir() string {
    if h := os.Getenv("HOME"); h != "" {
        return h
    }
    return os.Getenv("USERPROFILE") // windows
}

