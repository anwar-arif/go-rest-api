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
	flag.StringVar(&cfgPath, "config", "../test.config.yml", "config file path")
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

	// Initialize api client with timeout
	apiClient := &http.Client{Timeout: time.Minute * 2}

	/* ----TODO: uncomment this section and connect to db container and drop db in afterTest suit
	cfgMongo := config.GetMongo(cfgPath)
	cfgRedis := config.GetRedis(cfgPath)

	lgr := logger.DefaultOutStructLogger

	ctx := context.Background()

	// Initialize mongoDB
	mgo, err := infraMongo.New(ctx, cfgMongo.URL, cfgMongo.DBName, cfgMongo.DBTimeOut)
	Expect(err).NotTo(HaveOccurred())

	// initialize redis
	rds, err := infraRedis.New(ctx, cfgRedis.URL, cfgRedis.DbID, cfgRedis.DBTimeOut, lgr)
	Expect(err).NotTo(HaveOccurred())

	// initialize db
	db := infra.NewDB(mgo, rds)

	*/

	appBaseUrl := getAddressFromHostAndPort(cfgApp.Host, cfgApp.Port)

	var db *infra.DB

	framework.Root = framework.New(apiClient, cfgApp, db, appBaseUrl)

	/* -- TODO: uncomment to drop db if exist before test suit
	// drop db if exists
	By("dropping databases if exist")
	dbErr := framework.Root.DropDB(ctx)
	Expect(dbErr).NotTo(HaveOccurred())
	*/

})

var _ = AfterSuite(func() {
	//By("logout api test suite session")
	//framework.LogOut(framework.Root.Token)

	/* TODO: drop db after test suit
	ctx := context.Background()

	By("dropping database used for testing")
	err := framework.Root.DropDB(ctx)
	Expect(err).NotTo(HaveOccurred())
	By("dropped databases successfully")
	*/
})
