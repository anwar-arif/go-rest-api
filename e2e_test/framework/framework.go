package framework

import (
	"go-rest-api/config"
	"go-rest-api/infra/mongo"
	"net/http"
)

var (
	Root       *Framework
	SecretData *ConfidentialData
)

const (
	AuthenticationBearer = "bearer"
	AuthenticationToken  = "token"
	AuthenticationBoth   = "both"
	AuthenticationAny    = "any"
	AuthenticationNoAuth = "noAuth"
)

type ConfidentialData struct {
	UserName      string `json:"user_name"`
	Password      string `json:"password"`
	AuthBaseURL   string `json:"auth_base_url"`
	AuthSecretKey string `json:"auth_secret_key"`
}

type Framework struct {
	DB *mongo.Mongo
	//KV             *redis.Redis
	AppConfig      *config.Application
	ApiClient      *http.Client
	BaseURL, Token string
}

func New(
	apiClient *http.Client,
	appCfg *config.Application,
	db *mongo.Mongo,
	//kv *redis.Redis,
	baseURL string,
) *Framework {
	return &Framework{
		ApiClient: apiClient,
		AppConfig: appCfg,
		DB:        db,
		//KV:        kv,
		BaseURL: baseURL,
	}
}

func (f *Framework) DropDB() error {
	var err error
	//_ = f.KV.FlushDB()
	//err = f.DB.Database.Drop(context.TODO())
	//err = f.DB.Database().Collection().Drop(context.TODO())
	return err
}

func GetClient() *Framework {
	return Root
}
