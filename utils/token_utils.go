package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go-rest-api/model"
	"time"
)

type UserClaims struct {
	jwt.Claims
	Email string
	Role  string
}

func GenerateToken(user *model.User, secretKey string) (string, error) {
	appClaims := &UserClaims{
		Claims: jwt.MapClaims{
			"iat": jwt.NewNumericDate(time.Now()),
			"exp": jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
			"sub": user.UserName,
			"iss": "go-rest-api",
		},
		Email: user.Email,
		Role:  user.Role,
	}
	claims := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		appClaims,
	)
	secretInBytes := []byte(secretKey)

	token, err := claims.SignedString(secretInBytes)
	if err != nil {
		return "", err
	}

	return token, nil
}

func ValidateToken(tokenString string, secretKey string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true
	} else {
		fmt.Printf("invalid token %v\n", err.Error())
	}
	return false
}
