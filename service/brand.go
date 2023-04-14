package service

import (
	"context"
	"go-rest-api/api/response"
	"go-rest-api/logger"
	"go-rest-api/repo"
	"net/http"

	"go-rest-api/model"
	"go-rest-api/utils"
)

type BrandService interface {
	ListBrand(ctx context.Context, pager *utils.Pager) ([]model.BrandInfo, *response.Error)
	AddBrand(ctx context.Context, brand *model.BrandInfo) *response.Error
}

type brandService struct {
	log       logger.StructLogger
	brandRepo repo.BrandRepo
}

func NewBrandService(lgr logger.StructLogger, brandRepo repo.BrandRepo) BrandService {
	return &brandService{
		log:       lgr,
		brandRepo: brandRepo,
	}
}

// ListBrand ...
func (s *brandService) ListBrand(ctx context.Context, pager *utils.Pager) ([]model.BrandInfo, *response.Error) {
	fn := "ListBrands"
	tid := utils.GetTracingID(ctx)

	s.log.Println(fn, tid, "listing product brands from database")
	brands, err := s.brandRepo.ListBrands(ctx, "", pager.Skip, pager.Limit)
	if err != nil {
		s.log.Errorln(fn, tid, err.Error())
		return nil, response.NewError(http.StatusInternalServerError, err.Error())
	}

	s.log.Println(fn, tid, "sent response successfully")
	return brands, nil
}

// AddBrand ...
func (s *brandService) AddBrand(ctx context.Context, brand *model.BrandInfo) *response.Error {
	fn := "AddBrand"
	tid := utils.GetTracingID(ctx)

	s.log.Println(fn, tid, "inserting brands into database")
	err := s.brandRepo.Create(ctx, brand)
	if err != nil {
		s.log.Errorln(fn, tid, err.Error())
		return response.NewError(http.StatusInternalServerError, err.Error())
	}

	s.log.Println(fn, tid, "sent response successfully")
	return nil
}
