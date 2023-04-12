package service

import (
	"go-rest-api/config"
	"go-rest-api/infra"
	"go-rest-api/logger"
	"go-rest-api/repo"
)

// Service ...
type Service struct {
	log       logger.StructLogger
	brandRepo repo.BrandRepo
	userRepo  repo.UserRepo
}

func New(cfgDBTable *config.Table, db infra.DB, lgr logger.StructLogger) *Service {
	return &Service{
		log:       lgr,
		brandRepo: repo.NewBrand(cfgDBTable.BrandCollectionName, db),
		userRepo:  repo.NewUserRepo(cfgDBTable.UserCollectionName, db),
	}
}
