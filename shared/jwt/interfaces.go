package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"reflect"
	"sync"
)

//go:generate mockgen -source=./interfaces.go -destination=./interfaces.mock.gen.go -package=jwt

type Signer interface {
	CreateAccessToken(userId int64) (string, error)
	ParseWithClaims(token string) (claims *Claims, err error)
}

type rs256Signer struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func GetSigner() Signer {
	once.Do(func() {
		signer = &rs256Signer{
			publicKey:  readPublicKey(),
			privateKey: readPrivateKey(),
		}
	})

	return signer
}

var signer Signer
var once sync.Once
var basePath string

func init() { basePath, _ = os.Getwd() }

// TODO: need to change get key based on env if needed to go to production (separation for every ENV)

func readPrivateKey() *rsa.PrivateKey {
	buff, err := os.ReadFile(basePath + "/etc/jwt/private.key")
	if err != nil {
		panic("failed to read private key, err: " + err.Error())
	}

	block, _ := pem.Decode(buff)
	if block == nil {
		panic("failed to read pem block on private key")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic("failed to parse private key, err: " + err.Error())
	}

	return key
}

func readPublicKey() *rsa.PublicKey {
	buff, err := os.ReadFile(basePath + "/etc/jwt/public.key")
	if err != nil {
		panic("failed to read public key, err: " + err.Error())
	}

	block, _ := pem.Decode(buff)
	if block == nil {
		panic("failed to read pem block on public key")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic("failed to parse public key, err: " + err.Error())
	}

	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		panic("key is not rsa key, instead its " + reflect.TypeOf(key).Name())
	}

	return rsaKey
}
