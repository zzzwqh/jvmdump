package conf

import (
	"github.com/go-ini/ini"
)

type Config struct {
	Endpoint            string   `ini:"alioss.endpoint"`
	AccessKey           string   `ini:"alioss.accessKey"`
	AccessSecret        string   `ini:"alioss.accessSecret"`
	BucketName          string   `ini:"alioss.bucketName"`
	FolderName          string   `ini:"alioss.folderName"`
	DingtalkToken       string   `ini:"notify.dingtalkToken"`
	DingtalkSecret      string   `ini:"notify.dingtalkSecret"`
	DingtalkAt          []string `ini:"notify.dingtalkAt"`
	PrometheusAddr      string   `ini:"prometheus.address"`
	PrometheusThreshold float64  `ini:"prometheus.threshold"`
	RemoteDumpFileDir   string   `ini:"remote.dumpfile.dir"`
	LocalDumpFileDir    string   `ini:"local.dumpfile.dir"`
}

var inifile = "config.ini"
var GlobalConfig Config = Config{}

func init() {
	// 加载INI文件
	cfg, err := ini.Load(inifile)
	err = cfg.Section("").MapTo(&GlobalConfig)
	if err != nil {
		panic(err)
	}
}
