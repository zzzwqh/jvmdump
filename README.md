## 实现功能
通过 Prometheus 监控 Pod Jvm 指标，拉取到 Jvm Heap 堆内存水位
根据配置的阈值提前生成 dump 文件并上传到 OSS 
并通过钉钉 Webhook 发送告警

## 使用方法
克隆代码并构建，将配置文件和构建后的应用配置在相同目录下，可以写个 Dockerfile 和 Yaml 清单 

## 运行环境
需要在 Kubernetes 集群内部，赋予 rbac/rbac.yaml 权限，绑定 ServiceAccount pod-exec-read-logs-sa

## 相关配置
```bash
# 阿里云 OSS 配置
alioss.endpoint="xxxxxxxxxx"
alioss.accessKey="xxxxxxxxxxxx"
alioss.accessSecret="xxxxxxxxxxxxxxxxx"
alioss.bucketName="jvm-dumpfile"
alioss.folderName="dumpfile/test/"


# 钉钉群机器人 token
notify.dingtalkToken="xxxxxxxxxxxxxxxx"
notify.dingtalkSecret="xxxxxxxxxx"
notify.dingtalkAt="xxxxxxx"


# Prometheus 地址
prometheus.address="http://10.2.71.194:9090"

# jvm heap 水位阈值
# 0.01 表示当水位超过 1% 时就会生成 dump 文件并发送钉钉告警
prometheus.threshold=0.01



# 远程 jvm pod 生成 dump 文件的路径（如果不存在将会创建）
remote.dumpfile.dir="/opt/oom/"
# 本地 pod 存放 dump 文件的目录（如果不存在将会创建）
local.dumpfile.dir="/test/oom/"
```
