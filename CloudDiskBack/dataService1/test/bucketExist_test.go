package test

import (
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"testing"
)

func TestBucketExist(t *testing.T) {
	var obsClient, err = obs.New("OCGWMVEBEZ9ONDWM5OCB", "ukVqAjesme9oKJPSG7sNIkto7JR9l1HYl6PNs4Z9", "https://obs.cn-east-3.myhuaweicloud.com")
	if err != nil {
		t.Fatal("obsClient结构体创建失败")
	}
	_, err = obsClient.HeadBucket("data-storage101")
	if err == nil {
		fmt.Println("桶存在")
	} else {
		if obsError, ok := err.(obs.ObsError); ok {
			if obsError.StatusCode == 404 {
				fmt.Println("桶不存在")
			} else {
				fmt.Printf(err.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
	}
}
