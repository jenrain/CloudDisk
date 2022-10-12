package test

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"testing"
)

// 文件大小
//const chunkSize = 100 * 1024 * 1024 // 100MB
const chunkSize = 10 * 1024 * 1024 // 10MB

// 文件分片
func TestSeparateFile(t *testing.T) {
	fileInfo, err := os.Stat("videoTest.mp4")
	if err != nil {
		t.Fatal(err)
	}
	// 分片的个数
	chunkNum := int(math.Ceil(float64(fileInfo.Size()) / chunkSize))
	myFile, err := os.OpenFile("videoTest.mp4", os.O_RDONLY, 666)
	if err != nil {
		t.Fatal(err)
	}
	b := make([]byte, chunkSize)
	for i := 0; i < chunkNum; i++ {
		// 指定读取文件的起始位置
		myFile.Seek(int64(i*chunkSize), 0)
		if chunkSize > fileInfo.Size()-int64(i*chunkSize) {
			b = make([]byte, fileInfo.Size()-int64(i*chunkSize))
		}
		myFile.Read(b)
		f, err := os.OpenFile("./"+strconv.Itoa(i)+".chunk", os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
		f.Write(b)
		f.Close()
	}
	myFile.Close()
}

// 分片文件的合并
func TestMergeFile(t *testing.T) {
	myFile, err := os.OpenFile("videoTest2.mp4", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	fileInfo, err := os.Stat("videoTest.mp4")
	if err != nil {
		t.Fatal(err)
	}
	// 分片的个数
	chunkNum := int(math.Ceil(float64(fileInfo.Size()) / chunkSize))
	for i := 0; i < chunkNum; i++ {
		f, err := os.OpenFile("./"+strconv.Itoa(i)+".chunk", os.O_RDONLY, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			t.Fatal(err)
		}
		myFile.Write(b)
		f.Close()
	}
	myFile.Close()
}

// 文件一致性校验
func TestCheckFile(t *testing.T) {
	// 获取第一个文件的信息
	file1, err := os.OpenFile("videoTest.mp4", os.O_RDONLY, 0666)
	if err != nil {
		t.Fatal(err)
	}
	b1, err := ioutil.ReadAll(file1)
	if err != nil {
		t.Fatal(err)
	}
	// 获取第二个文件的信息
	// 获取第一个文件的信息
	file2, err := os.OpenFile("videoTest2.mp4", os.O_RDONLY, 0666)
	if err != nil {
		t.Fatal(err)
	}
	b2, err := ioutil.ReadAll(file2)
	if err != nil {
		t.Fatal(err)
	}
	// 比较哈希值
	s1 := fmt.Sprintf("%x", md5.Sum(b1))
	s2 := fmt.Sprintf("%x", md5.Sum(b2))
	fmt.Println(s1)
	fmt.Println(s2)
	fmt.Println(s1 == s2)
}
