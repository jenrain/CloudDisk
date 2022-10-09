package tools

import (
	"bytes"
	"core/define"
	"core/errorx"
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
func ObsUpload(r *http.Request) (string, error) {
	//u, _ := url.Parse(define.TencentCosBucket)
	//b := &cos.BaseURL{BucketURL: u}
	//c := cos.NewClient(b, &http.Client{
	//	Transport: &cos.AuthorizationTransport{
	//		SecretID:  define.TencentSecretID,  // 替换为用户的 SecretId，请登录访问管理控制台进行查看和管理，https://console.cloud.tencent.com/cam/capi
	//		SecretKey: define.TencentSecretKey, // 替换为用户的 SecretKey，请登录访问管理控制台进行查看和管理，https://console.cloud.tencent.com/cam/capi
	//	},
	//})
	//file, fileHeader, err := r.FormFile("file")
	//name := "cloud-disk/" + GetUUID() + path.Ext(fileHeader.Filename)
	//
	//buf := make([]byte, fileHeader.Size)
	//_, err = file.Read(buf)
	//_, err = c.Object.Put(context.Background(), name, strings.NewReader(string(buf)), nil)
	//
	////_, err = c.Object.Put(context.Background(), name, file, nil)
	//if err != nil {
	//	panic(err)
	//}
	//
	//// 返回资源的URL地址
	//return define.TencentCosBucket + "/" + name, nil
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
	key := "cloud-disk/" + GetUUID() + path.Ext(fileHeader.Filename)
	// 使用访问OBS
	input := &obs.PutObjectInput{}
	input.Bucket = "data-storage95"
	input.Key = key
	input.Body = strings.NewReader(string(buf))
	_, err = obsClient.PutObject(input)
	if err == nil {
		//fmt.Printf("RequestId:%s\n", output.RequestId)
		//fmt.Printf("ETag:%s\n", output.ETag)
	} else if obsError, ok := err.(obs.ObsError); ok {
		panic(obsError)
	}
	return define.HuaweiObsDawnLoadUrl + "/" + key, nil
}

// 创建ObsClient结构体
var obsClient, _ = obs.New(define.HuaweiObsAK, define.HuaweiObsSK, define.HuaweiObsEndPoint)

// ObsInitPartUpload 分片上传初始化
func ObsInitPartUpload(ext string) (string, string, error) {
	//u, _ := url.Parse(define.TencentCosBucket)
	//b := &cos.BaseURL{BucketURL: u}
	//c := cos.NewClient(b, &http.Client{
	//	Transport: &cos.AuthorizationTransport{
	//		SecretID:  define.TencentSecretID,  // 替换为用户的 SecretId，请登录访问管理控制台进行查看和管理，https://console.cloud.tencent.com/cam/capi
	//		SecretKey: define.TencentSecretKey, // 替换为用户的 SecretKey，请登录访问管理控制台进行查看和管理，https://console.cloud.tencent.com/cam/capi
	//	},
	//})
	//
	//name := "cloud-disk/" + GetUUID() + ext
	//// 可选opt,如果不是必要操作，建议上传文件时不要给单个文件设置权限，避免达到限制。若不设置默认继承桶的权限。
	//v, _, err := c.Object.InitiateMultipartUpload(context.Background(), name, nil)
	//if err != nil {
	//	return "", "", err
	//}
	//return name, v.UploadID, err
	inputInit := &obs.InitiateMultipartUploadInput{}
	inputInit.Bucket = define.HuaweiObsBucket
	inputInit.Key = define.HuaweiObsBucketRootFolder + "/" + GetUUID() + ext
	outputInit, err := obsClient.InitiateMultipartUpload(inputInit)
	return inputInit.Key, outputInit.UploadId, err
}

// ObsPartUpload 分片上传
func ObsPartUpload(r *http.Request) (string, error) {
	//u, _ := url.Parse(define.TencentCosBucket)
	//b := &cos.BaseURL{BucketURL: u}
	//c := cos.NewClient(b, &http.Client{
	//	Transport: &cos.AuthorizationTransport{
	//		SecretID:  define.TencentSecretID,  // 替换为用户的 SecretId，请登录访问管理控制台进行查看和管理，https://console.cloud.tencent.com/cam/capi
	//		SecretKey: define.TencentSecretKey, // 替换为用户的 SecretKey，请登录访问管理控制台进行查看和管理，https://console.cloud.tencent.com/cam/capi
	//	},
	//})
	//name := r.PostForm.Get("name")
	//UploadID := r.PostForm.Get("upload_id")
	//f, _, err := r.FormFile("file")
	//if err != nil {
	//	return "", err
	//}
	//
	//buf := bytes.NewBuffer(nil)
	//io.Copy(buf, f)
	//
	//// opt可选
	//partNumber, _ := strconv.Atoi(r.PostForm.Get("part_number"))
	//resp, err := c.Object.UploadPart(
	//	context.Background(), name, UploadID, partNumber, bytes.NewBuffer(buf.Bytes()), nil,
	//)
	//if err != nil {
	//	return "", err
	//}
	//return strings.Trim(resp.Header.Get("ETag"), "\""), nil
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
		Bucket:     define.HuaweiObsBucket,
		Key:        key,
		PartNumber: partNumber,
		UploadId:   uploadId,
		Body:       bytes.NewBuffer(buf.Bytes()),
	})
	return strings.Trim(outputUploadPart.ETag, "\""), err
	//fmt.Println("key: ", key)
	//fmt.Println("upload_id: ", uploadId)
	//fmt.Println("part_number: ", partNumber)
	//return "", nil
}

// ObsPartUploadComplete 分片上传完成
func ObsPartUploadComplete(key, uploadId string, obsObjects []obs.Part) (string, error) {
	//u, _ := url.Parse(define.TencentCosBucket)
	//b := &cos.BaseURL{BucketURL: u}
	//c := cos.NewClient(b, &http.Client{
	//	Transport: &cos.AuthorizationTransport{
	//		SecretID:  define.TencentSecretID,  // 替换为用户的 SecretId，请登录访问管理控制台进行查看和管理，https://console.cloud.tencent.com/cam/capi
	//		SecretKey: define.TencentSecretKey, // 替换为用户的 SecretKey，请登录访问管理控制台进行查看和管理，https://console.cloud.tencent.com/cam/capi
	//	},
	//})
	//opt := &cos.CompleteMultipartUploadOptions{}
	//opt.Parts = append(opt.Parts, cs...)
	//_, _, err := c.Object.CompleteMultipartUpload(
	//	context.Background(), name, uploadId, opt,
	//)
	//return err
	inputCompleteMultipart := &obs.CompleteMultipartUploadInput{}
	inputCompleteMultipart.Bucket = define.HuaweiObsBucket
	inputCompleteMultipart.Key = key
	inputCompleteMultipart.UploadId = uploadId

	inputCompleteMultipart.Parts = append(inputCompleteMultipart.Parts, obsObjects...)
	var err error
	_, err = obsClient.CompleteMultipartUpload(inputCompleteMultipart)
	return define.HuaweiObsDawnLoadUrl + "/" + key, err
}
