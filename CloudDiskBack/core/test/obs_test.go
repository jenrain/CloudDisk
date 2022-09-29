package test

import (
	"core/define"
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"testing"
)

func TestUploadFile(t *testing.T) {
	// 创建ObsClient结构体
	var obsClient, err = obs.New(define.HuaweiObsAK, define.HuaweiObsSK, define.HuaweiObsEndPoint)
	if err == nil {
		// 使用访问OBS
		input := &obs.PutFileInput{}
		input.Bucket = "data-storage95"
		input.Key = "cloud-disk/default.png"
		input.SourceFile = "../upload/default.png"
		output, err := obsClient.PutFile(input)

		if err == nil {
			fmt.Printf("RequestId:%s\n", output.RequestId)
			fmt.Printf("ETag:%s, StorageClass:%s\n", output.ETag, output.StorageClass)
		} else {
			if obsError, ok := err.(obs.ObsError); ok {
				fmt.Println(obsError.Code)
				fmt.Println(obsError.Message)
			} else {
				fmt.Println(err)
			}
		}
		// 关闭obsClient
		obsClient.Close()
	}
}
