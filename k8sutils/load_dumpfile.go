package k8sutils

import (
	"io"
	"jvmdump/conf"
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

// LoadDumpfile 远程到名称空间为 namespace 名字为 podName 的 containerName 容器中，将远程 pod 中的 dump 文件下载到本地 pod
func (k8sExecIns *K8sExec) LoadDumpfile(dumpfileName string) {
	// 指定要在目标 Pod 中执行的命令和参数
	command := []string{"tar", "cf", "-", conf.GlobalConfig.RemoteDumpFileDir + "/" + dumpfileName}

	// 创建 SPDY executor
	req := k8sExecIns.ClientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(k8sExecIns.PodName).
		Namespace(k8sExecIns.Namespace).
		SubResource("exec")
	req.VersionedParams(&corev1.PodExecOptions{
		Container: k8sExecIns.ContainerName,
		Command:   command,
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(k8sExecIns.RestConfig, "POST", req.URL())
	if err != nil {
		panic(err.Error())
	}

	// 创建本地存储路径
	localDir := conf.GlobalConfig.LocalDumpFileDir
	err = os.MkdirAll(localDir, 644)
	if err != nil {
		panic(err)
	}
	// 创建本地文件
	localFilePath := localDir + "/" + dumpfileName
	localFile, err := os.Create(localFilePath)
	if err != nil {
		panic(err.Error())
	}
	defer localFile.Close()

	// 将命令输出流和错误流重定向到本地文件
	err = executor.Stream(remotecommand.StreamOptions{
		Stdout: localFile,
		Stderr: os.Stderr,
	})
	if err != nil {
		panic(err.Error())
	}

	// 跳过 tar 文件头
	_, err = localFile.Seek(512, io.SeekStart)
	if err != nil {
		panic(err.Error())
	}
}
