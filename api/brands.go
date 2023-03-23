package api

import (
	"encoding/json"
	"net/http"

	"go-rest-api/api/response"
	"go-rest-api/logger"
	"go-rest-api/model"
	"go-rest-api/service"
	"go-rest-api/utils"
)

// BrandsController ...
type BrandsController struct {
	svc *service.Service
	lgr logger.StructLogger
}

// NewBrandsController ...
func NewBrandsController(svc *service.Service, lgr logger.StructLogger) *BrandsController {
	return &BrandsController{
		svc: svc,
		lgr: lgr,
	}
}

// ListBrand ...
func (cc *BrandsController) ListBrand(w http.ResponseWriter, r *http.Request) {
	tid := utils.GetTracingID(r.Context())
	pager, pageErr := utils.GetPager(r)
	if pageErr != nil {
		cc.lgr.Errorln("listBrands", tid, pageErr.Error())
		_ = response.ServeJSON(w, http.StatusBadRequest, nil, nil, pageErr.Error(), nil)
		return
	}

	cc.lgr.Println("listBrands", tid, "getting brands")
	result, err := cc.svc.ListBrand(r.Context(), pager)
	if err != nil {
		cc.lgr.Errorln("listBrands", tid, err.Error())
		_ = response.ServeJSON(w, err.StatusCode, nil, nil, err.Error(), nil)
		return
	}

	cc.lgr.Println("listBrands", tid, "sending response")
	_ = response.ServeJSON(w, http.StatusOK, pager.Prev, pager.Next, utils.SuccessMessage, result)
	return
}

// AddBrand ...
func (cc *BrandsController) AddBrand(w http.ResponseWriter, r *http.Request) {
	tid := utils.GetTracingID(r.Context())

	var b *model.BrandInfo
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		_ = response.ServeJSON(w, http.StatusBadRequest, nil, nil, utils.RequiredFieldMessage(), nil)
		return
	}

	cc.lgr.Println("AddBrand", tid, "inserting brand")
	err := cc.svc.AddBrand(r.Context(), b)
	if err != nil {
		cc.lgr.Errorln("AddBrand", tid, err.Error())
		_ = response.ServeJSON(w, err.StatusCode, nil, nil, err.Error(), nil)
		return
	}

	cc.lgr.Println("AddBrand", tid, "sending response")
	_ = response.ServeJSON(w, http.StatusOK, nil, nil, utils.SuccessMessage, nil)
	return
}