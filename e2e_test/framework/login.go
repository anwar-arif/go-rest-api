package framework

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
)

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginResp struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Data    LoginData `json:"data"`
}

func GetBearerToken(username, password string) string {
	login := &Login{Username: username, Password: password}

	payload, err := json.Marshal(login)
	Expect(err).NotTo(HaveOccurred())

	req, err := http.NewRequest("POST", SecretData.AuthBaseURL+"/login", bytes.NewBuffer(payload))
	By(fmt.Sprintf("user_name: %v, password: %v, url: %v", username, password, SecretData.AuthBaseURL+"/login"))
	Expect(err).NotTo(HaveOccurred())

	resp, err := Root.ApiClient.Do(req)
	Expect(err).NotTo(HaveOccurred())

	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).NotTo(HaveOccurred())
	By(fmt.Sprintf("received body from auth: %v", string(body)))
	err = resp.Body.Close()
	Expect(err).NotTo(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(200))

	var loginResp LoginResp
	err = json.Unmarshal(body, &loginResp)
	Expect(err).NotTo(HaveOccurred())

	return loginResp.Data.AccessToken
}

func LogOut(token string) {
	req, err := http.NewRequest("POST", SecretData.AuthBaseURL+"/logout", nil)
	Expect(err).NotTo(HaveOccurred())

	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := Root.ApiClient.Do(req)
	Expect(err).NotTo(HaveOccurred())

	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).NotTo(HaveOccurred())
	By(fmt.Sprintf("received body from auth: %v", string(body)))
	err = resp.Body.Close()
	Expect(err).NotTo(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(200))
}
