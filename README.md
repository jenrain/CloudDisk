# 基于go-zero实现的网盘系统

## 使用到的命令

**创建api服务**

```shell
goctl api new core
```

**使用api文件生成代码**

```shell
goctl api go -api core.api -dir . --style=goZero
```

**启动服务**

```shell
go run core.go -f etc/core-api.yaml
```

## 使用到的工具

**操作mysql**

Gorm：[https://github.com/go-gorm/gorm](https://github.com/go-gorm/gorm)

**操作Redis**

redigo：[https://github.com/gomodule/redigo](https://github.com/gomodule/redigo)

**邮箱验证**

jordan-wright：[https://github.com/jordan-wright/email](https://github.com/jordan-wright/email)

**生成uuid**

go.uuid：[https://github.com/satori/go.uuid](https://github.com/satori/go.uuid)

**生成jwt**

golang-jwt：[https://github.com/golang-jwt/jwt](https://github.com/golang-jwt/jwt)

**对象存储**

腾讯云：[https://github.com/tencentyun/cos-go-sdk-v5](https://github.com/tencentyun/cos-go-sdk-v5)
华为云：[https://github.com/huaweicloud/huaweicloud-sdk-go-obs](https://github.com/huaweicloud/huaweicloud-sdk-go-obs)