package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"strings"

	"go.uber.org/zap"
)

type Rsa struct {
	privateKey    string
	publicKey     string
	rsaPrivateKey *rsa.PrivateKey
	rsaPublicKey  *rsa.PublicKey
}

var ScRsa *Rsa

// 初始化rsa
func NewRsa(publicKey, privateKey string) *Rsa {
	ScRsa = &Rsa{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
	ScRsa.init()
	return ScRsa
}

func (thisRsa *Rsa) init() {
	if thisRsa.privateKey != "" {
		block, _ := pem.Decode([]byte(thisRsa.privateKey))
		//pkcs1
		if strings.Index(thisRsa.privateKey, "BEGIN RSA") > 0 {
			thisRsa.rsaPrivateKey, _ = x509.ParsePKCS1PrivateKey(block.Bytes)
		} else { //pkcs8
			privateKey, _ := x509.ParsePKCS8PrivateKey(block.Bytes)
			thisRsa.rsaPrivateKey = privateKey.(*rsa.PrivateKey)
		}
	}

	if thisRsa.publicKey != "" {
		block, _ := pem.Decode([]byte(thisRsa.publicKey))
		publickKey, _ := x509.ParsePKIXPublicKey(block.Bytes)
		thisRsa.rsaPublicKey = publickKey.(*rsa.PublicKey)
	}
	zap.L().Info("RSA init success...")
}

/**
 * 发布publicKey
 */
func (thisRsa *Rsa) GetPublicKey() string {
	return thisRsa.publicKey
}

// 生成公钥、私钥文件
func GenRsaKey(bits int) (private string, public string, err error) {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		zap.L().Error("GenRsaKey PrivateKey failed", zap.Error(err))
		return
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	private = string(pem.EncodeToMemory(block))

	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		zap.L().Error("GenRsaKey PublicKey failed", zap.Error(err))
		return
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	public = string(pem.EncodeToMemory(block))
	return
}
