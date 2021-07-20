package test

import (
	"go-rest-api/e2e_test/framework"
	"net/http"
)

type ApiResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SetAuthentication(req *http.Request, client *framework.Framework, authentication string) {
	switch authentication {
	case "bearer":
		req.Header.Set("Authorization", "Bearer "+client.Token)
	case "token":
		req.Header.Set("Secret-Key", framework.SecretData.AuthSecretKey)
	case "any":
	case "both":
		req.Header.Set("Authorization", "Bearer "+client.Token)
		req.Header.Set("Secret-Key", framework.SecretData.AuthSecretKey)
	}
}
