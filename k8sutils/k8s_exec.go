package k8sutils

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type K8sUtils interface {
	CreateDumpfile()
	LoadDumpfile()
}

// K8sExec 结构体，构造命令
type K8sExec struct {
	RestConfig    *rest.Config
	ClientSet     kubernetes.Interface
	PodName       string
	ContainerName string
	Namespace     string
}
