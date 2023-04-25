package k8sutils

import (
	"bytes"
	"fmt"
	"io"
	"jvmdump/conf"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

// CreateDumpfile 远程到名称空间为 namespace 名字为 podName 的 containerName 容器中，执行生成 dump 文件的命令
func (k8sExecIns *K8sExec) CreateDumpfile(dumpfileName string) {

	// oomDir 定义远程 Pod 中存放 dump 文件的位置
	oomDir := conf.GlobalConfig.RemoteDumpFileDir
	// 拼接生成 dumpfile 的绝对路径
	fullPath := oomDir + dumpfileName
	command := []string{"/bin/sh", "-c", "mkdir -p " + oomDir + " && if [ -f " + fullPath + " ]; then rm -rf " + fullPath + ";fi  && jmap -dump:live,format=b,file=" + fullPath + " 1"}

	// 创建 URL
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
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(k8sExecIns.RestConfig, "POST", req.URL())
	if err != nil {
		panic(err.Error())

	}

	var stdout, stderr io.Writer
	stdout = &bytes.Buffer{}
	stderr = &bytes.Buffer{}
	// 获取执行后的输出结果
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: stdout,
		Stderr: stderr,
		Tty:    false,
	})
	if err != nil {
		panic(err.Error())
	}
	// 打印输出结果
	fmt.Printf("stdout: %s\nstderr: %s\n", stdout, stderr)
}
