package utils

import (
	"go-rest-api/model"
	"testing"
)

var (
	user = model.User{
		UserID:   "",
		UserName: "Anwar35",
		Email:    "anwararif727@gmail.com",
		Role:     "admin",
	}
	apiSecretKey = "my_secret_key"
)

func TestGenerateToken(t *testing.T) {
	token, err := GenerateToken(&user, apiSecretKey)
	if err != nil {
		t.Errorf("err while generating token: %v", err.Error())
	}
	claims, err := GetClaimsFromToken(token, apiSecretKey)
	if err != nil {
		t.Errorf("can't parse token")
	}
	if (claims.Email != user.Email) || (claims.Role != user.Role) {
		t.Errorf("invalid token")
	}
}

func TestGetClaimsFromToken(t *testing.T) {
	token, err := GenerateToken(&user, apiSecretKey)
	if err != nil {
		t.Errorf("err while generating token: %v", err.Error())
	}
	claims, err := GetClaimsFromToken(token, apiSecretKey)
	if err != nil {
		t.Errorf("invalid token")
	}
	emailInToken := claims.Email
	roleInToken := claims.Role
	issuerInToken := claims.Issuer

	if (emailInToken != user.Email) || (roleInToken != user.Role) || (issuerInToken != Issuer) {
		t.Errorf("email or role in the token is tempered")
	}
}
