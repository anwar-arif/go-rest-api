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

func TestValidateToken(t *testing.T) {
	token, err := GenerateToken(&user, apiSecretKey)
	if err != nil {
		t.Errorf("err while generating token: %v", err.Error())
	}
	if !ValidateToken(token, apiSecretKey) {
		t.Errorf("invalid token")
	}
}
