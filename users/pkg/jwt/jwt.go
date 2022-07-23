package jwt

import (
	"errors"
	"fmt"
	gojwt "github.com/golang-jwt/jwt/v4"
	"strconv"
	"time"
)

var ErrInvalidToken = errors.New("invalid token")

type Jwt interface {
	Generate(userId uint, email string) (string, error)
	Validate(token string) (bool, error)
	Parse(token string) (userId uint, email string, err error)
}

type jwt struct {
	ttl      time.Duration
	baseTime time.Duration
	signKey  []byte
}

func NewJwtManager(tokenTtl time.Duration, baseTimeDelta time.Duration, signKey []byte) Jwt {
	if len(signKey) == 0 {
		panic("invalid key")
	}
	return &jwt{ttl: tokenTtl, baseTime: baseTimeDelta, signKey: signKey}
}

func (j *jwt) Generate(userId uint, email string) (string, error) {
	now := time.Now().Add(j.baseTime)
	claims := &gojwt.RegisteredClaims{
		ID:        strconv.Itoa(int(userId)),
		Issuer:    email,
		ExpiresAt: gojwt.NewNumericDate(now.Add(j.ttl)),
		NotBefore: gojwt.NewNumericDate(now),
	}
	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	return token.SignedString(j.signKey)
}

func (j *jwt) Validate(token string) (bool, error) {
	_, err := j.parse(token)
	if err != nil && err != ErrInvalidToken {
		return false, err
	}
	return err != ErrInvalidToken, nil
}

func (j *jwt) Parse(token string) (userId uint, email string, err error) {
	t, err := j.parse(token)
	if err != nil {
		return 0, "", err
	}
	claims := t.Claims.(*gojwt.RegisteredClaims)
	id, _ := strconv.Atoi(claims.ID)
	return uint(id), claims.Issuer, nil
}

func (j *jwt) parse(token string) (*gojwt.Token, error) {
	t, err := gojwt.ParseWithClaims(token, &gojwt.RegisteredClaims{}, func(token *gojwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*gojwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.signKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, ErrInvalidToken
	}
	return t, nil
}
