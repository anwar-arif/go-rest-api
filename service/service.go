package service

import (
	"go-rest-api/logger"
	"go-rest-api/repo"
)

// Service ...
type Service struct {
	log       logger.StructLogger
	brandRepo repo.BrandRepo
}

// New ...
func New(brandRepo repo.BrandRepo, lgr logger.StructLogger) *Service {
	return &Service{
		log:       lgr,
		brandRepo: brandRepo,
	}
}
