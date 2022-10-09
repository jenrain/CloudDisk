package define

import (
	"github.com/golang-jwt/jwt/v4"
)

type UserClaim struct {
	Id       int
	Identity string
	Name     string
	jwt.StandardClaims
}

var JwtKey = "cloud-disk-key"

// EmailPassword 邮箱密码
var EmailPassword = "RBQKKMWRKVAUISFS"

// CodeLength 验证码的长度
var CodeLength = 6

// TencentCosBucket 腾讯对象储存
var TencentCosBucket = "https://data-storage-1313761262.cos.ap-nanjing.myqcloud.com"
var TencentSecretKey = "5yc4JdnncGiz3xe9RZXHpCNNBoIgt9KS"
var TencentSecretID = "AKIDKhz72nfNruISRDqEOJ0d5JJ9nHbXJyPc"

// HuaweiObsEndPoint 华为云对象存储
var HuaweiObsEndPoint = "https://obs.cn-east-3.myhuaweicloud.com"
var HuaweiObsDawnLoadUrl = "https://data-storage95.obs.cn-east-3.myhuaweicloud.com"
var HuaweiObsBucket = "data-storage95"
var HuaweiObsBucketRootFolder = "cloud-disk"
var HuaweiObsAK = "OCGWMVEBEZ9ONDWM5OCB"
var HuaweiObsSK = "ukVqAjesme9oKJPSG7sNIkto7JR9l1HYl6PNs4Z9"

// PageSize 分页的默认参数
var PageSize = 20

// TokenExpireTime token的有效期
var TokenExpireTime int64 = 3600
var RefreshTokenExpireTime int64 = 2626560
