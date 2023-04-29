package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go-rest-api/e2e_test/framework"
	"go-rest-api/model"
	"go-rest-api/utils"
	"io"
	"net/http"
)

var _ = Describe("user api", func() {
	var f *framework.Framework
	signUpReq := model.SignUpRequest{
		UserName: "NightOwl",
		Email:    "myemail@gmail.com",
		Password: "my_password",
	}
	loginReq := model.LoginRequest{
		Email:    "myemail@gmail.com",
		Password: "my_password",
	}
	logOutReq := model.LogOutRequest{
		Email: "myemail@gmail.com",
	}
	userByEmailReq := model.GetUserByEmailRequest{
		Email: "myemail@gmail.com",
	}

	BeforeEach(func() {
		f = framework.GetFramework()
	})

	Context("sign up and login to get access token", func() {
		var accessToken string

		It("sign up should successful", func() {
			var buf bytes.Buffer
			url := f.BaseURL + "/api/v1/signup"
			buffErr := json.NewEncoder(&buf).Encode(signUpReq)
			Expect(buffErr).NotTo(HaveOccurred())

			req, reqErr := http.NewRequest(http.MethodPost, url, &buf)
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
			url := f.BaseURL + "/api/v1/login"

			buffErr := json.NewEncoder(&buf).Encode(loginReq)
			Expect(buffErr).NotTo(HaveOccurred())

			req, reqErr := http.NewRequest(http.MethodPost, url, &buf)
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

		It("should fetch logged in user data", func() {
			var buf bytes.Buffer
			url := f.BaseURL + "/api/v1/users"
			buffErr := json.NewEncoder(&buf).Encode(userByEmailReq)
			Expect(buffErr).NotTo(HaveOccurred())

			req, reqErr := http.NewRequest(http.MethodGet, url, &buf)
			Expect(reqErr).NotTo(HaveOccurred())

			req.Header.Set(utils.AuthorizationKey, "Bearer "+accessToken)

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

			userResponse := model.GetUserByEmailResponse{
				Email: respMap["email"].(string),
			}
			Expect(userByEmailReq.Email).To(Equal(userResponse.Email))
			By("get user by email successful")
		})

		It("should logout successful", func() {
			var buf bytes.Buffer
			url := f.BaseURL + "/api/v1/logout"
			buffErr := json.NewEncoder(&buf).Encode(logOutReq)
			Expect(buffErr).NotTo(HaveOccurred())

			req, reqErr := http.NewRequest(http.MethodPost, url, &buf)
			Expect(reqErr).NotTo(HaveOccurred())

			req.Header.Set(utils.AuthorizationKey, "Bearer "+accessToken)

			resp, respErr := f.ApiClient.Do(req)
			Expect(respErr).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			body, parseErr := io.ReadAll(resp.Body)
			Expect(parseErr).NotTo(HaveOccurred())

			var res ApiResponse
			err := json.Unmarshal(body, &res)
			Expect(err).NotTo(HaveOccurred())

			Expect(res.Message).To(Equal(utils.LogoutSuccessful))
			By("logout api successful")

		})
	})
})
