package utils

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"go-rest-api/api/response"
	"go-rest-api/model"
	"time"
)

type ClaimsKey int

var claimsKey ClaimsKey

func SetJWTClaimsContext(ctx context.Context, claims UserClaims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

func JWTClaimsFromContext(ctx context.Context) (*UserClaims, bool) {
	claims, ok := ctx.Value(claimsKey).(*UserClaims)
	return claims, ok
}

type UserClaims struct {
	*jwt.RegisteredClaims
	Email string `json:"email"`
	Role  string `json:"role"`
}

var (
	Issuer   = "https://anwararif.com"
	Audience = "go-rest-api"
)

func GenerateToken(user *model.User, secretKey string) (string, error) {
	appClaims := &UserClaims{
		&jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Millisecond * -1)),
			Subject:   user.UserName,
			Issuer:    Issuer,
			Audience:  []string{Audience},
		},
		user.Email,
		user.Role,
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

func GetClaimsFromToken(tokenString string, secretKey string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		if _, ok := token.Claims.(*UserClaims); !ok && !token.Valid {
			return nil, response.TokenExpired
		}
		issuer, err := token.Claims.GetIssuer()
		if (err != nil) || (issuer != Issuer) {
			return nil, jwt.ErrTokenInvalidIssuer
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
