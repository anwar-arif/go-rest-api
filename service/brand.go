package service

import (
	"context"
	"go-rest-api/api/response"
	"net/http"

	"go-rest-api/model"
	"go-rest-api/utils"
)

// ListBrand ...
func (s *Service) ListBrand(ctx context.Context, pager *utils.Pager) ([]model.BrandInfo, *response.Error) {
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
func (s *Service) AddBrand(ctx context.Context, brand *model.BrandInfo) *response.Error {
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
