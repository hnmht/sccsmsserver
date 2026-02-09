package jwt

import (
	"sccsmsserver/i18n"
	"sccsmsserver/pub"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
)

type MyClaims struct {
	UserID   int32  `json:"userid"`
	UserCode string `json:"usercode"`
	jwt.StandardClaims
}

// Generate token
func GenToken(userID int32, usercode string, tokenID string) (tokenString string, expireTime int64, err error) {
	expireTime = time.Now().Add(pub.TokenExpireDuration).Unix()
	c := MyClaims{
		UserID:   userID,
		UserCode: usercode,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime,
			Issuer:    pub.TokenIssuer,
			Id:        tokenID,
		},
	}
	// create signature object
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// signed token
	tokenString, err = token.SignedString(pub.TokenSecret)
	if err != nil {
		zap.L().Error("GenToken token.SignedString failed:", zap.Error(err))
	}
	return
}

// Parse Token
func ParseToken(tokenString string) (*MyClaims, i18n.ResKey) {
	var mc = new(MyClaims)

	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		return pub.TokenSecret, nil
	})

	if err != nil {
		return nil, i18n.CodeInvalidToken
	}

	if token.Valid {
		return mc, i18n.StatusOK
	}

	return nil, i18n.CodeInvalidToken
}
