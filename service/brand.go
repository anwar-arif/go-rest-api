package service

import (
	"context"
	"go-rest-api/api/response"
	"net/http"

	"go-rest-api/logger"
	"go-rest-api/model"
	"go-rest-api/utils"
)

// SetLogger ...
func (c *Service) SetLogger(l logger.StructLogger) {
	c.log = l
}

// ListBrand ...
func (c *Service) ListBrand(ctx context.Context, pager *utils.Pager) ([]model.BrandInfo, *response.Error) {
	tid := utils.GetTracingID(ctx)

	c.log.Println("ListBrands", tid, "listing product brands from database")
	brands, err := c.brandRepo.ListBrands(ctx, "", pager.Skip, pager.Limit)
	if err != nil {
		c.log.Errorln("ListBrands", tid, err.Error())
		return nil, response.NewError(http.StatusInternalServerError, err.Error())
	}

	c.log.Println("ListBrands", tid, "sent response successfully")
	return brands, nil
}

// AddBrand ...
func (c *Service) AddBrand(ctx context.Context, brand *model.BrandInfo) *response.Error {
	tid := utils.GetTracingID(ctx)

	c.log.Println("AddBrand", tid, "inserting brands into database")
	err := c.brandRepo.Create(ctx, brand)
	if err != nil {
		c.log.Errorln("AddBrand", tid, err.Error())
		return response.NewError(http.StatusInternalServerError, err.Error())
	}

	c.log.Println("AddBrand", tid, "sent response successfully")
	return nil
}
