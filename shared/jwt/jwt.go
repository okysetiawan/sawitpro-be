package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	UserId int64
}

func (t *rs256Signer) CreateAccessToken(userId int64) (string, error) {

	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "sawitpro",
			ExpiresAt: jwt.NewNumericDate(time.Now()),
		},
		UserId: userId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(t.privateKey)
}

func (t *rs256Signer) ParseWithClaims(token string) (claims *Claims, err error) {
	claims = new(Claims)

	_, err = jwt.ParseWithClaims(token, claims, func(_ *jwt.Token) (interface{}, error) { return t.publicKey, nil })
	if err != nil {
		return nil, err
	}

	return claims, nil
}
