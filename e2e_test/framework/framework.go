package framework

import (
	"context"
	"errors"
	"go-rest-api/config"
	infra "go-rest-api/infra/db"
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
	AuthBaseURL   string `json:"auth_base_url"`
	AuthSecretKey string `json:"auth_secret_key"`
}

type Framework struct {
	DB        *infra.DB
	AppConfig *config.Application
	ApiClient *http.Client
	BaseURL   string
	Token     string
}

func New(
	apiClient *http.Client,
	appCfg *config.Application,
	db *infra.DB,
	baseURL string,
) *Framework {
	return &Framework{
		ApiClient: apiClient,
		AppConfig: appCfg,
		DB:        db,
		BaseURL:   baseURL,
	}
}

func (f *Framework) DropDB(ctx context.Context) error {
	redisErr := f.DB.Redis.FlushDB(ctx)
	mgoErr := f.DB.Mongo.DropDB(ctx)
	err := errors.Join(redisErr, mgoErr)
	return err
}

func GetFramework() *Framework {
	return Root
}
