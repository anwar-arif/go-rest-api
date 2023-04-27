package e2e_test

import (
	"context"
	"flag"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"go-rest-api/config"
	"go-rest-api/e2e_test/framework"
	_ "go-rest-api/e2e_test/test"
	infraMongo "go-rest-api/infra/mongo"
	"net/http"
	"strconv"
	"testing"
	"time"
)

// Ginkgo generated this function to kick off our unit tests. Hook into it to define factories.

func TestSignatures(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Api Test Suite")
}

var (
	cfgPath string
)

func init() {
	flag.StringVar(&cfgPath, "config", "test.config.yaml", "config file path")
}

func getAddressFromHostAndPort(host string, port int) string {
	addr := host
	if host == "" {
		addr = "http://localhost"
	}
	if port != 0 {
		addr = addr + ":" + strconv.Itoa(port)
	}
	return addr
}

var _ = BeforeSuite(func() {
	By("going for api, db, kv initialization")
	// get configuration
	cfgApp := config.GetApp(cfgPath)
	cfgMongo := config.GetMongo(cfgPath)

	// Initialize api client with timeout
	apiClient := &http.Client{Timeout: time.Minute * 2}
	ctx := context.Background()

	// Initialize mongoDB and  client
	db, err := infraMongo.New(ctx, cfgMongo.URL, cfgMongo.DBName, cfgMongo.DBTimeOut)
	Expect(err).NotTo(HaveOccurred())

	// Initialize redis and client
	//kv := infraRedis.New(cfgPath, "test")

	appBaseUrl := getAddressFromHostAndPort(cfgApp.Host, cfgApp.Port)
	framework.SecretData = &framework.ConfidentialData{
		UserName:      viper.GetString("secret_data.user_name"),
		Password:      viper.GetString("secret_data.password"),
		AuthBaseURL:   viper.GetString("app.auth_base_url") + "/api/v1/auth",
		AuthSecretKey: viper.GetString("app.api_secret_key"),
	}

	Expect(err).NotTo(HaveOccurred())
	framework.Root = framework.New(apiClient, cfgApp, db, appBaseUrl)
	By("going for login attempt")

	token := framework.GetBearerToken(framework.SecretData.UserName, framework.SecretData.Password)
	framework.Root.Token = token
})

var _ = AfterSuite(func() {
	By("logout api test suite session")
	framework.LogOut(framework.Root.Token)

	err := framework.Root.DropDB()
	Expect(err).NotTo(HaveOccurred())
})
