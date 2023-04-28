package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go-rest-api/e2e_test/framework"
	"go-rest-api/model"
	"io"
	"net/http"
)

var _ = Describe("user api", func() {
	var f *framework.Framework
	signUpReq := model.SignUpRequest{
		UserName: "NightOwl",
		Email:    "anwararif727@gmail.com",
		Password: "my_password",
	}
	loginReq := model.LoginRequest{
		Email:    "anwararif727@gmail.com",
		Password: "my_password",
	}

	BeforeEach(func() {
		f = framework.GetFramework()
	})

	Context("sign up and login to get access token", func() {
		var accessToken string

		It("sign up should successful", func() {
			var buf bytes.Buffer
			signUpUrl := f.BaseURL + "/api/v1/signup"
			buffErr := json.NewEncoder(&buf).Encode(signUpReq)
			Expect(buffErr).NotTo(HaveOccurred())

			req, reqErr := http.NewRequest(http.MethodPost, signUpUrl, &buf)
			Expect(reqErr).NotTo(HaveOccurred())

			By("sending sign up request")
			resp, respErr := f.ApiClient.Do(req)
			Expect(respErr).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusCreated))

			body, parseErr := io.ReadAll(resp.Body)
			Expect(parseErr).NotTo(HaveOccurred())
			By(fmt.Sprintf("sing up response %v", string(body)))

			var expected = &ApiResponse{
				Message: "Created",
				Data: map[string]interface{}{
					"user_name": signUpReq.UserName,
					"email":     signUpReq.Email,
				},
			}

			var res ApiResponse
			err := json.Unmarshal(body, &res)
			Expect(err).NotTo(HaveOccurred())

			Expect(res.Message).To(Equal(expected.Message))
			Expect(res.Data).To(Equal(expected.Data))

			By("sign up success")
		})

		It("login should successful", func() {
			var buf bytes.Buffer
			loginUrl := f.BaseURL + "/api/v1/login"

			buffErr := json.NewEncoder(&buf).Encode(loginReq)
			Expect(buffErr).NotTo(HaveOccurred())

			req, reqErr := http.NewRequest(http.MethodPost, loginUrl, &buf)
			Expect(reqErr).NotTo(HaveOccurred())

			By("sending login request")
			resp, respErr := f.ApiClient.Do(req)
			Expect(respErr).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			body, parseErr := io.ReadAll(resp.Body)
			Expect(parseErr).NotTo(HaveOccurred())

			var res ApiResponse
			err := json.Unmarshal(body, &res)
			Expect(err).NotTo(HaveOccurred())

			respMap, ok := res.Data.(map[string]interface{})
			Expect(ok).To(Equal(true))

			loginRespData := &model.LoginResponse{
				Email:       respMap["email"].(string),
				AccessToken: respMap["access_token"].(string),
			}

			accessToken = loginRespData.AccessToken

			By(fmt.Sprintf("received response %v", res.Data))
			Expect(loginRespData.Email).To(Equal(loginReq.Email))
			Expect(accessToken).NotTo(BeEmpty())

			By("login successful")
		})
	})
})
