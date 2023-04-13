package api

import (
	"encoding/json"
	"go-rest-api/api/response"
	"go-rest-api/infra"
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
	fn := "CreateUser"
	tid := utils.GetTracingID(r.Context())

	var u *model.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		uc.lgr.Errorln(fn, tid, err.Error())
		_ = response.Serve(w, http.StatusBadRequest, utils.RequiredFieldMessage(), nil)
		return
	}

	// user exists with same email
	existingUser, err := uc.svc.GetByEmail(r.Context(), &u.Email)
	if err.Error() != infra.ErrNotFound.Error() {
		uc.lgr.Errorln(fn, tid, err.Error())
		_ = response.Serve(w, http.StatusInternalServerError, "failed to create user", nil)
		return
	}
	if existingUser != nil {
		uc.lgr.Errorln(fn, tid, "this email is already in use")
		_ = response.Serve(w, http.StatusConflict, "email already in use", nil)
		return
	}

	if err := uc.svc.CreateUser(r.Context(), u); err != nil {
		uc.lgr.Errorln(fn, tid, err.Error())
		_ = response.Serve(w, err.StatusCode, err.Error(), nil)
		return
	}

	_ = response.Serve(w, http.StatusOK, utils.SuccessMessage, &model.GetUserByEmailResponse{
		UserName: u.UserName,
		Email:    u.Email,
	})
	return
}

func (uc *UsersController) GetByEmail(w http.ResponseWriter, r *http.Request) {
	fn := "GetByEmail"
	tid := utils.GetTracingID(r.Context())

	var getByEmailReq *model.GetUserByEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&getByEmailReq); err != nil {
		_ = response.Serve(w, http.StatusBadRequest, utils.RequiredFieldMessage("email"), nil)
		return
	}

	user, err := uc.svc.GetByEmail(r.Context(), &getByEmailReq.Email)
	if err != nil {
		uc.lgr.Errorln(fn, tid, err.Error())
		_ = response.Serve(w, err.StatusCode, err.Error(), nil)
		return
	}

	_ = response.Serve(w, http.StatusOK, utils.SuccessMessage, user)
	return
}
