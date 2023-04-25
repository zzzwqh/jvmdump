package ossutils

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"jvmdump/conf"
	"os"
)

// UploadDumpfile 传入文件路径 , 文件名称
func UploadDumpfile(filePath, filename string) {
	// yourEndpoint填写Bucket对应的Endpoint，以华东1（杭州）为例，填写为https://ossutils-cn-hangzhou.aliyuncs.com。其它Region请按实际情况填写。
	client, err := oss.New(conf.GlobalConfig.Endpoint, conf.GlobalConfig.AccessKey, conf.GlobalConfig.AccessSecret)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	bucketName := conf.GlobalConfig.BucketName

	// 填写存储空间名称，例如 examplebucket。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 依次填写Object的完整路径（例如 exampledir/exampleobject.txt ）和本地文件的完整路径（例如D:\\localpath\\examplefile.txt）。
	err = bucket.PutObjectFromFile(conf.GlobalConfig.FolderName+filename, filePath+"/"+filename)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

}
