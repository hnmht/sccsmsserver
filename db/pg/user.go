package pg

import (
	"crypto/md5"
	"encoding/hex"
	"sccsmsserver/pub"
)

// encryptPassword 加密密码
func EncryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(pub.Secret))

	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}
