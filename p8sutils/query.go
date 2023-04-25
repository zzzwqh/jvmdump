package p8sutils

import (
	"context"
	"encoding/json"
	"fmt"
	"jvmdump/conf"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
)

// Data Json 反序列化要转成的结构体
type Data struct {
	PodName       string `json:"pod,omitempty"`
	PodPort       string `json:"port,omitempty"`
	Application   string `json:"application,omitempty"`
	Area          string `json:"area,omitempty"`
	ID            string `json:"id,omitempty"`
	Instance      string `json:"instance,omitempty"`
	Job           string `json:"job,omitempty"`
	K8sNamespace  string `json:"namespace,omitempty"`
	ContainerName string `json:"container,omitempty"`
}

// ResultList 包含了两个列表，一个是从 prometheus 获取的 Json 列表，另外一个是 Metrics 指标列表
type ResultList struct {
	JsonData    []Data
	MetricsData []string
}

const (
	// promSql jvm heap 状态查询语句
	promSql = "jvm_memory_used_bytes{area=\"heap\", job!=\"springboot-metrics-v1.0\", job!=\"nacos-cluster\"} / jvm_memory_max_bytes{area=\"heap\", job!=\"springboot-metrics-v1.0\", job!=\"nacos-cluster\"}"
)

// PromSqlQuery 执行查询，返回超出阈值的指标数组
func PromSqlQuery() (resultList ResultList) {
	// 创建 prometheus api 客户端
	client, err := api.NewClient(api.Config{
		Address: conf.GlobalConfig.PrometheusAddr,
	})
	if err != nil {
		panic(err)
	}

	// 创建 prometheus api 查询器
	queryAPI := v1.NewAPI(client)
	// 将 prometheus sql 语句和阈值放到一起
	fullSql := promSql + ">" + strconv.FormatFloat(conf.GlobalConfig.PrometheusThreshold, 'E', -1, 64)
	// 查询指标，获取 result 结果
	result, _, err := queryAPI.Query(context.Background(), fullSql, time.Now())
	if err != nil {
		panic(err)
	}
	// 将获取到的结果分割
	podList := strings.Split(result.String(), "\n")
	// 定义接收 Json 反序列化后的结构体数组
	var dataSlice []Data
	// 定义接收 metrics 指标的数组
	var metricsSlice []string
	for _, item := range podList {
		// 获取指标的 Json 字符串（根据 => 做分隔符，取到数组第一个元素）
		parts := strings.Split(item, "=>")
		if len(parts) != 2 {
			fmt.Println("获取 Json 数据出错 ~")

		}
		dataStr := strings.TrimSpace(parts[0])
		// 一些狗屎代码，将查询到的数据，转换成 Json 格式（将 = 符号换成 : 符号，并将 Key 值加上了 "" 双引号）,后面需要研究下 prometheus 官方提供的方法
		dataStr = strings.ReplaceAll(dataStr, "=", "\":")
		dataStr = strings.ReplaceAll(dataStr, "{", "{\"")
		dataStr = strings.ReplaceAll(dataStr, ", ", ",\"")
		var data Data
		err = json.Unmarshal([]byte(dataStr), &data)
		if err != nil {
			fmt.Println(err.Error(), "unable to parse json data")
		}
		// 获取的没问题，就将实例结构体传入到 dataSlice
		dataSlice = append(dataSlice, data)
		// 获取当前 jvm 百分比查询指标
		metricsParts := strings.Split(strings.TrimSpace(parts[1]), "@[")
		if len(metricsParts) != 2 {
			fmt.Println("invalid timestamp format")
			return
		}
		metric := strings.TrimSpace(metricsParts[0])
		// 获取的没问题，就将实例结构体传入到 metricsSlice
		metricsSlice = append(metricsSlice, metric)
	}
	// 返回 Json 数据反序化列后的对象数组，和 Metrics 指标数组，封装到 resultList 结构体中
	resultList = ResultList{JsonData: dataSlice, MetricsData: metricsSlice}
	return resultList
}
