package main

import (
	"context"
	"fmt"
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp" // add this to avoid "no Auth Provider found for name "gcp""
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

type k8s struct {
	clientset kubernetes.Interface
}

func newK8s() (*k8s, error){
	path := os.Getenv("HOME") + "/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", path)
	if err != nil {
		return nil, err
	}

	client := k8s{}
	client.clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &client, nil
}

func (o *k8s) getVersion() (string, error) {
	version, err := o.clientset.Discovery().ServerVersion();
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", version), nil
}

func (o *k8s) printNamespaces() ([]string, error) {
	nlist, err := o.clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var res []string
	for _, ns := range nlist.Items {
		res = append(res, ns.Name)
	}
	return res, nil
}

func (o *k8s) printPodsInfo() {
	podlist, err := o.clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, pod := range podlist.Items {
		fmt.Printf("Namespace %s, PodName %s, PodLabels %s\n", pod.Namespace, pod.Name, pod.Labels)
	}
}
func main() {
	k8s, err := newK8s()
	if err != nil {
		fmt.Println(err)
		return
	}
	version, err := k8s.getVersion()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(version)
	if ns, err := k8s.printNamespaces(); err != nil {
		panic(err.Error())
	} else {
		fmt.Println(ns)
	}

	k8s.printPodsInfo()


}

