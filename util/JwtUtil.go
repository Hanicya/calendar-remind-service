package util

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

var jwtKey []byte = []byte{87, 50, 212, 174, 174, 20, 169, 128, 174, 65, 243, 200, 149, 127, 66, 190, 184, 59, 49, 124, 151, 214, 51, 166, 184, 193, 9, 175, 69, 40, 205, 38}

func GenerateJwt(key any, method jwt.SigningMethod, claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(method, claims)
	return token.SignedString(key)
}

func ParseJwt(key any, jwtStr string, options ...jwt.ParserOption) (int, error) {
	token, err := jwt.Parse(jwtStr, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	}, options...)
	if err != nil {
		return -1, err
	}
	// 校验 Claims 对象是否有效
	if !token.Valid {
		return -1, errors.New("invalid token")
	}
	issuerStr, _ := token.Claims.GetIssuer()
	issuerInt, err := strconv.Atoi(issuerStr)
	if err != nil {
		return -1, err // 如果转换失败，返回错误
	}
	return issuerInt, nil
}

func main() {
	//生成32字节（256位）的密钥 jwtKey := make([]byte, 32)
	if _, err := rand.Read(jwtKey); err != nil {
		panic(err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "1",
		"sub": "JiNiTaiMei",
		"aud": "iKun",
		"exp": time.Now().Add(time.Second * 10).UnixMilli(),
	})
	jwtStr, err := token.SignedString(jwtKey)
	if err != nil {
		panic(err)
	}
	fmt.Println(jwtStr)

	// 解析 jwt
	claims, err := ParseJwt(jwtKey, jwtStr, jwt.WithExpirationRequired())
	if err != nil {
		panic(err)
	}
	fmt.Println(claims)

}
