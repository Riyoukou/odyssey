package utils

import (
	"context"
	"fmt"
	"log"

	rolloutsv1beta1 "github.com/openkruise/kruise-rollout-api/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func CreateKubernetesClientset(kubernetesConfig string) *kubernetes.Clientset {
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubernetesConfig))
	if err != nil {
		log.Fatalf("Failed to create REST config: %v", err)
	}

	// 创建 Clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func CreateKruiseClientset(kubernetesConfig string) *rolloutsv1beta1.Clientset {
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubernetesConfig))
	if err != nil {
		log.Fatalf("Failed to create REST config: %v", err)
	}

	// 创建 Clientset
	clientset, err := rolloutsv1beta1.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func RolloutGetOptions(clientset *dynamic.DynamicClient, namespace, appName string) *unstructured.Unstructured {
	// 定义 Rollout 资源的 GroupVersionResource
	rolloutResource := schema.GroupVersionResource{
		Group:    "rollouts.kruise.io",
		Version:  "v1beta1",
		Resource: "rollouts",
	}

	// 获取指定的 Rollout 资源
	rollout, err := clientset.Resource(rolloutResource).Namespace(namespace).Get(context.Background(), appName, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("获取 Rollout 资源失败: %v\n", err)
		fmt.Printf("namespace: %s, appName: %s\n", namespace, appName)
		return nil
	}
	return rollout
}

func RolloutUpdateOptions(clientset *dynamic.DynamicClient, namespace, appName string, rollout *unstructured.Unstructured) error {
	// 定义 Rollout 资源的 GroupVersionResource
	rolloutResource := schema.GroupVersionResource{
		Group:    "rollouts.kruise.io",
		Version:  "v1beta1",
		Resource: "rollouts",
	}

	_, err := clientset.Resource(rolloutResource).Namespace(namespace).UpdateStatus(context.Background(), rollout, metav1.UpdateOptions{})
	if err != nil {
		fmt.Printf("Rollout 更新失败: %v\n", err)
		return err
	}

	return nil
}
