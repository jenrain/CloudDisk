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
# 启动rabbitmq
docker run -d --hostname rabbit-svr --name rabbit -p 5672:5672 -p 15672:15672 -p 25672:25672 -v /opt/module/rabbitmq/data:/var/lib/rabbitmq rabbitmq:management

# 启动canal
cd /opt/module/canal/bin/
./startup.sh

# 启动前端命令
npm run serve

# 启动API层的服务
go run apiService.go -f etc/core-api.yaml
# 启动Data层的服务
go run dataService.go -f etc/dataservice-api.yaml
# 后端服务可以直接双击start.bat启动
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