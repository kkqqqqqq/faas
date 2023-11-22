package essemble

import (
	"context"
	"flag"
	"fmt"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/clientcmd"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/homedir"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

func main() {
	var kubeconfigPath string
	flag.StringVar(&kubeconfigPath, "kubeconfig", "", "Path to the kubeconfig file (default is $HOME/.kube/config)")
	flag.Parse()

	if kubeconfigPath == "" {
		// 如果未提供 kubeconfig 文件路径，则默认使用 $HOME/.kube/config
		homeDir := homedir.HomeDir()
		kubeconfigPath = homeDir + "/.kube/config"
	}
	// 加载 kubeconfig 文件，以获取集群配置
	config,
		err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		fmt.Printf("Error loading kubeconfig: %v\n", err)
		panic(err.Error())
	}
	// 创建上下文和取消函数
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 创建 Kubernetes 客户端
	//clientset, err := kubernetes.NewForConfig(config)
	//if err != nil {
	//	panic(err.Error())
	//}

	// 获取K8s集群节点名称 for _, node := range nodes.Items
	//nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})

	// 创建 Metrics 客户端
	metricsClient, err := metricsv.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating Metrics client: %v\n", err)
		panic(err.Error())
	}

	// 定义要获取指标的节点名称
	nodeName := "node_name" // 替换为实际的节点名称

	// 定义获取指标的时间间隔
	interval := 30 * time.Second

	// 创建周期性获取节点资源指标的函数
	wait.Until(func() {
		nodeMetrics, err := metricsClient.MetricsV1beta1().NodeMetricses().Get(ctx, nodeName, metav1.GetOptions{})
		if err != nil {
			fmt.Printf("Error getting metrics for node %s: %v\n", nodeName, err)
			return
		}

		cpuUsage := nodeMetrics.Usage.Cpu()

		memoryUsage := nodeMetrics.Usage.Memory()

		fmt.Printf("Node: %s\n", nodeName)
		fmt.Printf("CPU Usage: %s\n", cpuUsage)
		fmt.Printf("Memory Usage: %s\n", memoryUsage)
	}, interval, ctx.Done())

	// 阻塞主线程以保持运行
	select {}
}

func getCPU(model string) (cpu float32) {
	return Models[model].cpu
}

func getMemory(model string) (memory int) {
	return Models[model].memory
}
