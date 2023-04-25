package notify

import (
	"github.com/CodyGuo/dingtalk"
	"github.com/CodyGuo/dingtalk/pkg/robot"
	"github.com/CodyGuo/glog"
	"jvmdump/conf"
	"strings"

	"io/ioutil"
)

// SendMsgDingtalk 根据传入的 bucketUrl,namespace,podName,ipAddress,application,metricsRes 生成钉钉 markdown 消息
func SendMsgDingtalk(bucketUrl string, namespace string, podName string, ipAddress string, application string, metricsRes string) {

	glog.SetFlags(glog.LglogFlags)
	webHook := "https://oapi.dingtalk.com/robot/send?access_token=" + conf.GlobalConfig.DingtalkToken
	// 机器人安全设置页面，加签一栏勾选之后下面显示的 SEC 开头的字符串
	secret := conf.GlobalConfig.DingtalkSecret
	dt := dingtalk.New(webHook, dingtalk.WithSecret(secret))

	// markdown 类型
	atPersons := strings.Join(conf.GlobalConfig.DingtalkAt, " ")

	markdownTitle := "JVM Heap Warnning"
	markdownText := "#### **JVM Heap 告警**  \n" +
		"#### **Namespace:** " + namespace + " \n" +
		"#### **PodName:** " + podName + " \n" +
		"#### **PodIP:** " + ipAddress + " \n" +
		"#### **Application:** " + application + " \n" +
		"#### **<font color=\"red\">JVM 堆内存水位: " + metricsRes + "% </font>** \n" +
		atPersons + "\n" +
		"> 已生成 dump 文件到 OSS [点击下载](" +
		bucketUrl +
		")\n"
	atMobiles := robot.SendWithAtMobiles(conf.GlobalConfig.DingtalkAt)
	if err := dt.RobotSendMarkdown(markdownTitle, markdownText, atMobiles); err != nil {
		glog.Fatal(err)
	}
	printResult(dt)

}

func printResult(dt *dingtalk.DingTalk) {
	response, err := dt.GetResponse()
	if err != nil {
		glog.Fatal(err)
	}
	reqBody, err := response.Request.GetBody()
	if err != nil {
		glog.Fatal(err)
	}
	reqData, err := ioutil.ReadAll(reqBody)
	if err != nil {
		glog.Fatal(err)
	}
	glog.Infof("发送消息成功, message: %s", reqData)
}
