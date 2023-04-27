package utils

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"go-rest-api/api/response"
	"go-rest-api/model"
	"net/http"
	"strings"
	"time"
)

type ClaimsKey int

var claimsKey ClaimsKey

func SetJWTClaimsContext(ctx context.Context, requestId string, claims UserClaims) context.Context {
	fmt.Println("requestId: ", requestId, ", claimsKey: ", claimsKey)
	return context.WithValue(ctx, requestId, claims)
}

func JWTClaimsFromContext(ctx context.Context, requestId string) (*UserClaims, bool) {
	fmt.Println("requestId: ", requestId, ", claimsKey: ", claimsKey)
	claims, ok := ctx.Value(requestId).(*UserClaims)
	return claims, ok
}

type UserClaims struct {
	*jwt.RegisteredClaims
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

var (
	Issuer   = "https://anwararif.com"
	Audience = "go-rest-api"
)

func GenerateToken(user *model.User, secretKey string) (string, error) {
	appClaims := &UserClaims{
		&jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
			// TODO: use expiresIn param to calculate ExpiresAt
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 300)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Millisecond * -1)),
			Subject:   user.UserName,
			Issuer:    Issuer,
			Audience:  []string{Audience},
		},
		user.UserID,
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

func ClaimsFromRequest(r *http.Request) (*UserClaims, error) {
	bearer := r.Header.Get(AuthorizationKey)
	token := strings.TrimPrefix(bearer, "Bearer ")
	secretKey := viper.GetString("app.api_secret_key")
	return GetClaimsFromToken(token, secretKey)
}
