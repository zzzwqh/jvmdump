package main

import (
	"fmt"
	"jvmdump/conf"
	"jvmdump/k8sutils"
	"jvmdump/notify"
	"jvmdump/ossutils"
	"jvmdump/p8sutils"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"strconv"
	"strings"
)

func main() {
	var dumpFilename string
	var jvmHeapPercent string

	// 创建连接 kubernetes 的客户端
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// 通过 prometheus 查询，获取超出阈值的 instance 列表
	resultList := p8sutils.PromSqlQuery()
	// 循环取出 instance 列表中的元素
	for index, item := range resultList.JsonData {
		// 根据 index 获取 resultList.MetricsData 当前指标的值
		jvmHeapPercent = resultList.MetricsData[index]
		fmt.Println(jvmHeapPercent)
		floatRes, err := strconv.ParseFloat(jvmHeapPercent, 64)
		if err != nil {
			fmt.Println(err.Error())
		}
		// 将查询到的指标，转换成百分比数值
		floatRes = floatRes * 100
		metricsResult := fmt.Sprintf("%.2f", floatRes)

		// 根据 instance 的指标，生成 dumpfile 文件名称，以作区分，因为可能重复生成，给文件加上时间戳，但缺点是会累积 dumpfile，需要清理
		// dumpFilename = item.Application + "-" + strings.Split(item.Instance, ":")[0] + "-oom-" + time.Now().Format("20060102150405") + ".dump"
		// 根据 instance 的指标，生成 dumpfile 文件名称，以作区分，如果重复生成会删除并重新创建（检查和删除功能在 CreateDumpfile 中使用 bash 实现）
		dumpFilename = item.Application + "-" + strings.Split(item.Instance, ":")[0] + "-oom.dump"
		fmt.Println(dumpFilename)

		// 生成 k8sExecIns 结构体对象，定义 config/client/pod 等信息
		var k8sExecIns = k8sutils.K8sExec{
			RestConfig:    config,
			ClientSet:     clientSet,
			PodName:       item.PodName,
			ContainerName: item.ContainerName,
			Namespace:     item.K8sNamespace,
		}
		// 通过调用 kubernetes API 执行生成 dump 文件
		k8sExecIns.CreateDumpfile(dumpFilename)

		// 通过调用 kubernetes API 执行拷贝 dump 文件到本地 pod
		k8sExecIns.LoadDumpfile(dumpFilename)

		// 传送 dumpfile 到 OSS
		ossutils.UploadDumpfile(conf.GlobalConfig.LocalDumpFileDir, dumpFilename)

		// 拼接 OSS 存储中的文件链接
		fullFilePath := "https://" + conf.GlobalConfig.BucketName + "." + conf.GlobalConfig.Endpoint + "/" + dumpFilename

		// 发送钉钉消息,传入参数
		notify.SendMsgDingtalk(fullFilePath, item.K8sNamespace, item.PodName, strings.Split(item.Instance, ":")[0], item.Application, metricsResult)
	}
}
