package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
)

type MyClaims struct {
	UserID   int32  `json:"userid"`
	UserCode string `json:"usercode"`
	jwt.StandardClaims
}

var ErrorInvalidToken = errors.New("无效Token")

const TokenExpireDuration = time.Hour * 2

var mySercet = []byte("这是一个加密密码")

// GenToken 生成JWT
func GenToken(userID int32, usercode string, tokenID string) (tokenString string, expireTime int64, err error) {
	expireTime = time.Now().Add(TokenExpireDuration).Unix()
	//解析token
	c := MyClaims{
		UserID:   userID,
		UserCode: usercode,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime, // 过期时间
			Issuer:    "seacloud", //签发人
			Id:        tokenID,    //token ID
		},
	}
	//使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	//转换
	tokenString, err = token.SignedString(mySercet)

	if err != nil {
		zap.L().Error("GenToken token.SignedString failed:", zap.Error(err))
	}
	return
}

// ParseToken 解析jwt
func ParseToken(tokenString string) (*MyClaims, error) {
	//解析token
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		return mySercet, nil
	})

	if err != nil {
		return nil, err
	}

	if token.Valid {
		return mc, nil
	}

	return nil, ErrorInvalidToken
}
