package tools

import (
	"bytes"
	"core/define"
	"core/errorx"
	"core/models"
	"core/mq"
	"crypto/md5"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/jordan-wright/email"
	uuid "github.com/satori/go.uuid"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"path"
	"strconv"
	"strings"
	"time"
)

// StringToMD5 将字符串做MD5哈希
func StringToMD5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// FileToMD5 将文件做MD5哈希
func FileToMD5(file multipart.File, size int64) (string, error) {
	var err error
	b := make([]byte, size)
	_, err = file.Read(b)
	return StringToMD5(string(b)), err
}

func GenerateToken(id int, identity string, name string, second int64) (string, error) {
	uc := define.UserClaim{
		Id:       id,
		Identity: identity,
		Name:     name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(second)).Unix(),
		},
	}
	// 获取到token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	// 将token加密
	signedString, err := token.SignedString([]byte(define.JwtKey))
	if err != nil {
		return "", err
	}
	return signedString, nil
}

// ParseToken 解析token
func ParseToken(tokenString string) (*define.UserClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &define.UserClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(define.JwtKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*define.UserClaim); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// MailSendCode 邮箱验证码发送
func MailSendCode(mail, code string) error {
	e := email.NewEmail()
	e.From = "Cloud Disk <onemorelight234@163.com>"
	e.To = []string{"501124524@qq.com"}
	e.Subject = "验证码发送测试"
	e.HTML = []byte("您的验证码为：<h1>" + code + "</h1>")
	err := e.SendWithTLS("smtp.163.com:465", smtp.PlainAuth("", "onemorelight234@163.com", define.EmailPassword, "smtp.163.com"),
		&tls.Config{
			InsecureSkipVerify: true,
			ServerName:         "smtp.163.com",
		})
	if err != nil {
		return err
	}
	return nil
}

// RandCode 生成验证码
func RandCode() string {
	s := "1234567890"
	code := ""
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < define.CodeLength; i++ {
		code += string(s[rand.Intn(len(s))])
	}
	return code
}

// GetUUID 获取UUID
func GetUUID() string {
	return uuid.NewV4().String()
}

// ObsUpload 上传文件到Cos
func ObsUpload(r *http.Request, nodes []models.Obs, fileIdentity string) (string, error) {
	// 创建ObsClient结构体
	var obsClient, err = obs.New(define.HuaweiObsAK, define.HuaweiObsSK, define.HuaweiObsEndPoint)
	defer obsClient.Close()
	if err != nil {
		return "", errorx.NewDefaultError("创建ObsClient结构体失败")
	}
	// 解析文件
	file, fileHeader, err := r.FormFile("file")
	buf := make([]byte, fileHeader.Size)
	_, err = file.Read(buf)
	key := "cloud-disk/" + fileIdentity + path.Ext(fileHeader.Filename)
	// 统计文件存入了哪些节点中
	nodeList := ""
	// 分别向随机生成的节点发送文件
	for _, node := range nodes {
		nodeList += node.HuaweiObsBucket + ","
		fmt.Println("正在向节点：", node.HuaweiObsBucket, "上传数据")
		input := &obs.PutObjectInput{}
		input.Bucket = node.HuaweiObsBucket
		input.Key = key
		input.Body = strings.NewReader(string(buf))
		_, err = obsClient.PutObject(input)
		if err == nil {
		} else if obsError, ok := err.(obs.ObsError); ok {
			panic(obsError)
		}
	}
	return nodeList, nil
}

// ObsInitPartUpload 分片上传初始化
func ObsInitPartUpload(ext, fileIdentity string, node models.Obs) (string, string, error) {
	obsClient, _ := obs.New(node.HuaweiObsAK, node.HuaweiObsSK, node.HuaweiObsEndPoint)
	inputInit := &obs.InitiateMultipartUploadInput{}
	inputInit.Bucket = node.HuaweiObsBucket
	inputInit.Key = define.HuaweiObsBucketRootFolder + "/" + fileIdentity + ext
	outputInit, err := obsClient.InitiateMultipartUpload(inputInit)
	return inputInit.Key, outputInit.UploadId, err
}

// ObsPartUpload 分片上传
func ObsPartUpload(r *http.Request, node models.Obs) (string, error) {
	obsClient, _ := obs.New(node.HuaweiObsAK, node.HuaweiObsSK, node.HuaweiObsEndPoint)
	key := r.PostForm.Get("key")
	uploadId := r.PostForm.Get("upload_id")
	f, _, err := r.FormFile("file")
	if err != nil {
		return "", errorx.NewDefaultError(err.Error())
	}
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, f)
	partNumber, _ := strconv.Atoi(r.PostForm.Get("part_number"))
	var outputUploadPart *obs.UploadPartOutput
	outputUploadPart, err = obsClient.UploadPart(&obs.UploadPartInput{
		Bucket:     node.HuaweiObsBucket,
		Key:        key,
		PartNumber: partNumber,
		UploadId:   uploadId,
		Body:       bytes.NewBuffer(buf.Bytes()),
	})
	return strings.Trim(outputUploadPart.ETag, "\""), err
}

// ObsPartUploadComplete 分片上传完成
func ObsPartUploadComplete(key, uploadId string, obsObjects []obs.Part, node models.Obs) (string, error) {
	obsClient, _ := obs.New(node.HuaweiObsAK, node.HuaweiObsSK, node.HuaweiObsEndPoint)
	inputCompleteMultipart := &obs.CompleteMultipartUploadInput{}
	inputCompleteMultipart.Bucket = node.HuaweiObsBucket
	inputCompleteMultipart.Key = key
	inputCompleteMultipart.UploadId = uploadId

	inputCompleteMultipart.Parts = append(inputCompleteMultipart.Parts, obsObjects...)
	var err error
	_, err = obsClient.CompleteMultipartUpload(inputCompleteMultipart)
	return node.HuaweiObsDawnLoadUrl + "/" + key, err
}

// ObsPartUploadTheRestPart 分片上传继续上传剩下的部分
func ObsPartUploadTheRestPart(key, path string, node models.Obs, nodes []models.Obs) (err error) {
	// 先将文件下载到本地
	err = obsDownLoadFile(key, path, node)
	if err != nil {
		fmt.Println("文件下载错误")
	}
	// 开始上传剩下的部分
	obsClient, err := obs.New(node.HuaweiObsAK, node.HuaweiObsSK, node.HuaweiObsEndPoint)
	for _, s := range nodes {
		fmt.Println("开始向剩下的：", s.HuaweiObsBucket, "节点上传文件")
		input := &obs.PutFileInput{}
		input.Bucket = s.HuaweiObsBucket
		input.Key = key
		input.SourceFile = ".\\download\\" + key
		_, err = obsClient.PutFile(input)
	}
	return err
}

// ObsDownLoadFile 断点续传下载
func obsDownLoadFile(key, path string, node models.Obs) error {
	obsClient, err := obs.New(node.HuaweiObsAK, node.HuaweiObsSK, node.HuaweiObsEndPoint)
	input := &obs.DownloadFileInput{}
	input.Bucket = node.HuaweiObsBucket
	input.Key = key
	input.DownloadFile = ".\\download\\" + key // localfile为下载对象的本地文件全路径
	input.EnableCheckpoint = true              // 开启断点续传模式
	input.PartSize = 9 * 1024 * 1024           // 指定分段大小为9MB
	input.TaskNum = 5                          // 指定分段下载时的最大并发数
	_, err = obsClient.DownloadFile(input)
	return err
}

// GetFilePath 获取文件的path
func GetFilePath(uuid, ext string) (path string, err error) {
	// 获取所有节点列表
	nodes := mq.GetDataServers()
	input := &obs.GetObjectMetadataInput{}
	for _, node := range nodes {
		obsClient, _ := obs.New(node.HuaweiObsAK, node.HuaweiObsSK, node.HuaweiObsEndPoint)
		input.Bucket = node.HuaweiObsBucket
		input.Key = node.HuaweiObsBucketRootFolder + "/" + uuid + ext
		_, err = obsClient.GetObjectMetadata(input)
		if err == nil {
			path = node.HuaweiObsDawnLoadUrl + "/" + input.Key
			break
		}
	}
	return
}
