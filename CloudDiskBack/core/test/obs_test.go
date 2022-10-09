package test

import (
	"core/define"
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"strconv"
	"testing"
)

// 创建ObsClient结构体
var obsClient, _ = obs.New(define.HuaweiObsAK, define.HuaweiObsSK, define.HuaweiObsEndPoint)

// 简单上传
func TestUploadFileToOBS(t *testing.T) {
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
}

// 初始化分片上传
func TestInitPartUploadToOBS(t *testing.T) {
	var uploadId = ""
	// 初始化分段上传任务
	inputInit := &obs.InitiateMultipartUploadInput{}
	inputInit.Bucket = "data-storage95"
	inputInit.Key = "cloud-disk/videoTest.mp4"
	outputInit, err := obsClient.InitiateMultipartUpload(inputInit)
	if err != nil {
		if obsError, ok := err.(obs.ObsError); ok {
			fmt.Println(obsError.Code)
			fmt.Println(obsError.Message)
		} else {
			fmt.Println(err)
		}
	}
	uploadId = outputInit.UploadId
	fmt.Printf("UploadId:%s\n", uploadId)
}

// 分片上传
func TestPartUploadToOBS(t *testing.T) {
	etagSlice := make([]string, 0)
	var err error
	for i := 0; i < 14; i++ {
		var outputUploadPart *obs.UploadPartOutput
		outputUploadPart, err = obsClient.UploadPart(&obs.UploadPartInput{
			Bucket:     "data-storage95",
			Key:        "cloud-disk/videoTest.mp4",
			PartNumber: i + 1,
			UploadId:   "000001838A2F4DC79016308B91399BBC",
			SourceFile: "./" + strconv.Itoa(i) + ".chunk",
		})
		fmt.Println("ETag: ", outputUploadPart.ETag)
		etagSlice = append(etagSlice, outputUploadPart.ETag)
	}
	if err != nil {
		if obsError, ok := err.(obs.ObsError); ok {
			fmt.Println(obsError.Code)
			fmt.Println(obsError.Message)
		} else {
			fmt.Println(err)
		}
	}
	fmt.Println(etagSlice)
}

// 分片上传完成
func TestPartUploadToOBSComplete(t *testing.T) {
	uploadId := "000001838A2F4DC79016308B91399BBC"
	inputCompleteMultipart := &obs.CompleteMultipartUploadInput{}
	inputCompleteMultipart.Bucket = "data-storage95"
	inputCompleteMultipart.Key = "cloud-disk/videoTest.mp4"
	inputCompleteMultipart.UploadId = uploadId
	//for i := 0; i < 14; i++ {
	//	inputCompleteMultipart.Parts = append(inputCompleteMultipart.Parts, obs.Part{
	//		PartNumber:   i + 1,
	//		ETag:         "",
	//	})
	//}
	inputCompleteMultipart.Parts = []obs.Part{
		obs.Part{PartNumber: 1, ETag: "2c34c52a859ed209b1d32fc9bd5f761f"},
		obs.Part{PartNumber: 2, ETag: "3e8b5685a4ff6476d4e4559dfd758094"},
		obs.Part{PartNumber: 3, ETag: "ea0c241c1df977202f3382f6fa2c6499"},
		obs.Part{PartNumber: 4, ETag: "575a3d555e7fcc5aa8be41354d7de562"},
		obs.Part{PartNumber: 5, ETag: "67cc3734762f38a202afe34b6113481d"},
		obs.Part{PartNumber: 6, ETag: "cca91f0c8d5ce90ae9e5b48050590425"},
		obs.Part{PartNumber: 7, ETag: "5f8ee1af29e0b00421f791b274416939"},
		obs.Part{PartNumber: 8, ETag: "f4a7ff83e7c5fee9f415dbc0a8e12176"},
		obs.Part{PartNumber: 9, ETag: "bb39075b52af9f6b2adf97d1d4c854ad"},
		obs.Part{PartNumber: 10, ETag: "8292c596c39023170c84ffe376ab8997"},
		obs.Part{PartNumber: 11, ETag: "5970408c19d49a63bce7796af5740687"},
		obs.Part{PartNumber: 12, ETag: "0e4fea5915b5959d544e9b643567a35a"},
		obs.Part{PartNumber: 13, ETag: "eade8418a1f1b98cf42c76682d1de8df"},
		obs.Part{PartNumber: 14, ETag: "30b52f63b0f1e49432c4c05051c5e590"},
	}
	outputCompleteMultipart, err := obsClient.CompleteMultipartUpload(inputCompleteMultipart)
	if err != nil {
		if obsError, ok := err.(obs.ObsError); ok {
			fmt.Println(obsError.Code)
			fmt.Println(obsError.Message)
		} else {
			fmt.Println(err)
		}
	}
	fmt.Printf("RequestId:%s\n", outputCompleteMultipart.RequestId)
	fmt.Printf("Location:%s, Bucket:%s, Key:%s, ETag:%s\n", outputCompleteMultipart.Location, outputCompleteMultipart.Bucket, outputCompleteMultipart.Key, outputCompleteMultipart.ETag)

}
