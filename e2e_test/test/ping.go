package test

import (
	"encoding/json"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go-rest-api/e2e_test/framework"
	"io/ioutil"
	"net/http"
)

var (
	client *framework.Framework
)

var _ = Describe("Ping Api", func() {
	BeforeEach(func() {
		client = framework.GetClient()
	})

	Context("Calling Ping", func() {
		It("should return a ping response", func() {
			req, err := http.NewRequest("GET", client.BaseURL, nil)
			Expect(err).NotTo(HaveOccurred())

			By(fmt.Sprintf("set %v authentication to header", framework.AuthenticationBearer))
			SetAuthentication(req, client, framework.AuthenticationBearer)

			By("sending ping request")
			resp, err := client.ApiClient.Do(req)
			Expect(err).NotTo(HaveOccurred())

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			By(fmt.Sprintf("received body from ping request: %v", string(body)))

			// status should be created 200 for this api
			By("checking ping request status code")
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var expected = &ApiResponse{
				Message: "success",
				Data:    map[string]interface{}{"method": "GET", "service_name": "Go rest api"},
			}
			// check response body
			var res ApiResponse
			err = json.Unmarshal(body, &res)
			Expect(res.Message).To(Equal(expected.Message))
			Expect(res.Data).To(Equal(expected.Data))
			By("Ping request succeeded")
		})
	})
})
