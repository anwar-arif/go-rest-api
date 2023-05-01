package e2e_test

import (
	"flag"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go-rest-api/config"
	"go-rest-api/e2e_test/framework"
	_ "go-rest-api/e2e_test/test"
	infra "go-rest-api/infra/db"
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
	flag.StringVar(&cfgPath, "config", "../test.config.yaml", "config file path")
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
	//cfgMongo := config.GetMongo(cfgPath)
	//cfgRedis := config.GetRedis(cfgPath)
	//
	//lgr := logger.DefaultOutStructLogger

	// Initialize api client with timeout
	apiClient := &http.Client{Timeout: time.Minute * 2}
	//ctx := context.Background()

	// Initialize mongoDB
	//mgo, err := infraMongo.New(ctx, cfgMongo.URL, cfgMongo.DBName, cfgMongo.DBTimeOut)
	//Expect(err).NotTo(HaveOccurred())
	//
	//// initialize redis
	//rds, err := infraRedis.New(ctx, cfgRedis.URL, cfgRedis.DbID, cfgRedis.DBTimeOut, lgr)
	//Expect(err).NotTo(HaveOccurred())

	// initialize db
	//db := infra.NewDB(mgo, rds)
	var db *infra.DB
	//Expect(err).NotTo(HaveOccurred())

	appBaseUrl := getAddressFromHostAndPort(cfgApp.Host, cfgApp.Port)
	//framework.SecretData = &framework.ConfidentialData{
	//	AuthBaseURL:   viper.GetString("app.auth_base_url") + "/api/v1/auth",
	//	AuthSecretKey: viper.GetString("app.api_secret_key"),
	//}

	//Expect(err).NotTo(HaveOccurred())
	framework.Root = framework.New(apiClient, cfgApp, db, appBaseUrl)

	// drop db if exists
	//By("dropping databases if exist")
	//dbErr := framework.Root.DropDB(ctx)
	//Expect(dbErr).NotTo(HaveOccurred())
	//By("going for login attempt")

	//token := framework.GetBearerToken(framework.SecretData.UserName, framework.SecretData.Password)
	//framework.Root.Token = token
})

var _ = AfterSuite(func() {
	//By("logout api test suite session")
	//framework.LogOut(framework.Root.Token)

	//ctx := context.Background()

	//By("dropping database used for testing")
	//err := framework.Root.DropDB(ctx)
	//Expect(err).NotTo(HaveOccurred())
	//By("dropped databases successfully")
})
