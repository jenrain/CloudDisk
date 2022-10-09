package test

import (
	"core/define"
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"testing"
)

// 流式下载
func TestStreamDownload(t *testing.T) {
	var obsClient, err = obs.New(define.HuaweiObsAK, define.HuaweiObsSK, define.HuaweiObsEndPoint)
	if err == nil {
		input := &obs.GetObjectInput{}
		input.Bucket = "data-storage95"
		input.Key = "cloud-disk/681b51ae-590d-4de0-b204-d1235b62974b.jpg"
		output, err := obsClient.GetObject(input)
		if err == nil {
			defer output.Body.Close()
			fmt.Printf("StorageClass:%s, ETag:%s, ContentType:%s, ContentLength:%d, LastModified:%s\n",
				output.StorageClass, output.ETag, output.ContentType, output.ContentLength, output.LastModified)
			p := make([]byte, 1024)
			var readErr error
			var readCount int
			// 读取对象内容
			for {
				readCount, readErr = output.Body.Read(p)
				if readCount > 0 {
					fmt.Printf("%s", p[:readCount])
				}
				if readErr != nil {
					break
				}
			}
		} else if obsError, ok := err.(obs.ObsError); ok {
			fmt.Printf("Code:%s\n", obsError.Code)
			fmt.Printf("Message:%s\n", obsError.Message)
		}
		// 关闭obsClient
		obsClient.Close()
	}
}
