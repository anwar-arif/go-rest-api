package api

import (
	"encoding/json"
	"go-rest-api/api/response"
	"go-rest-api/logger"
	"go-rest-api/model"
	"go-rest-api/service"
	"go-rest-api/utils"
	"net/http"
)

type UsersController struct {
	svc *service.Service
	lgr logger.StructLogger
}

func NewUsersController(svc *service.Service, lgr logger.StructLogger) *UsersController {
	return &UsersController{
		svc: svc,
		lgr: lgr,
	}
}

func (uc *UsersController) CreateUser(w http.ResponseWriter, r *http.Request) {
	tid := utils.GetTracingID(r.Context())

	var u *model.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		uc.lgr.Errorln("CreateUser", tid, err.Error())
		_ = response.ServeJSON(w, http.StatusBadRequest, nil, nil, utils.RequiredFieldMessage(), nil)
		return
	}

	if err := uc.svc.CreateUser(r.Context(), u); err != nil {
		uc.lgr.Errorln("CreateUser", tid, err.Error())
		_ = response.ServeJSON(w, err.StatusCode, nil, nil, err.Error(), nil)
		return
	}

	_ = response.ServeJSON(w, http.StatusOK, nil, nil, utils.SuccessMessage, nil)
	return
}

func (uc *UsersController) GetByEmail(w http.ResponseWriter, r *http.Request) {
	tid := utils.GetTracingID(r.Context())

	var email *string
	if err := json.NewDecoder(r.Body).Decode(&email); err != nil {
		_ = response.ServeJSON(w, http.StatusBadRequest, nil, nil, utils.RequiredFieldMessage("email"), nil)
		return
	}

	user, err := uc.svc.GetByEmail(r.Context(), email)
	if err != nil {
		uc.lgr.Errorln("GetByEmail", tid, err.Error())
		_ = response.ServeJSON(w, err.StatusCode, nil, nil, err.Error(), nil)
		return
	}

	_ = response.ServeJSON(w, http.StatusOK, nil, nil, utils.SuccessMessage, user)
	return
}
