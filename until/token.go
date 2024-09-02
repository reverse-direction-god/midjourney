package until

import (
	"fmt"
	"strconv"

	"github.com/golang-jwt/jwt"
)

var jwtSecret = []byte("ylsy")

func Encryption(userID uint) string {
	// 创建一个新的token对象，指定签名方法和声明
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
	})

	// 使用secret key签名并获得完整的编码token
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return ""
	}

	return tokenString
}
func Decrypt(tokenString string) string {
	// 解析token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证token的签名方法是否为预期的方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return ""
	}

	// 验证token并获取声明
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["userID"].(float64)
		userId := int(userID)
		userid := strconv.Itoa(userId)
		return userid
	} else {
		return ""
	}
}
