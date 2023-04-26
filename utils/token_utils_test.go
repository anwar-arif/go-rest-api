package utils

import (
	"fmt"
	"go-rest-api/model"
	"testing"
)

var (
	user = model.User{
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
	fmt.Println("token: ", token)
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
	emailInToken := claims["email"].(string)
	roleInToken := claims["role"].(string)

	if (emailInToken != user.Email) || (roleInToken != user.Role) {
		t.Errorf("email or role in the token is tempered")
	}
}
