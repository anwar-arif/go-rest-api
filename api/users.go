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
		uc.lgr.Errorf(fn, tid, "error while parsing user payload: %v\n", err.Error())
		_ = response.Serve(w, http.StatusBadRequest, utils.RequiredFieldMessage(), nil)
		return
	}

	// user exists with same email
	existingUser, err := uc.svc.GetUserByEmail(r.Context(), &u.Email)
	if (err != nil) && (err.Error() != infra.ErrNotFound.Error()) {
		uc.lgr.Errorf(fn, tid, "error while fetching user: %v\n", err.Error())
		_ = response.Serve(w, http.StatusInternalServerError, "failed to create user", nil)
		return
	}

	if existingUser != nil {
		uc.lgr.Errorln(fn, tid, "this email is already in use")
		_ = response.Serve(w, http.StatusConflict, "email already in use", nil)
		return
	}

	salt := utils.GenerateSalt()
	if salt == nil {
		uc.lgr.Errorln(fn, tid, "failed to generate salt")
		_ = response.Serve(w, http.StatusInternalServerError, "failed to create user", nil)
		return
	}

	u.Salt = *salt
	u.Password = utils.HashPassword(u.Password, u.Salt)

	if err := uc.svc.CreateUser(r.Context(), u); err != nil {
		uc.lgr.Errorln(fn, tid, err.Error())
		_ = response.Serve(w, err.StatusCode, err.Error(), nil)
		return
	}

	_ = response.Serve(w, http.StatusCreated, http.StatusText(http.StatusCreated), &model.GetUserByEmailResponse{
		UserName: u.UserName,
		Email:    u.Email,
	})
	return
}

func (uc *UsersController) LogIn(w http.ResponseWriter, r *http.Request) {
	fn := "LogIn"
	tid := utils.GetTracingID(r.Context())

	var loginReq *model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		_ = response.Serve(w, http.StatusBadRequest, response.CannotProcessRequest, nil)
		return
	}

	user, err := uc.svc.GetAuthUserByEmail(r.Context(), &loginReq.Email)
	if (err != nil) && (err.Error() != infra.ErrNotFound.Error()) {
		uc.lgr.Errorln(fn, tid, http.StatusText(http.StatusInternalServerError))
		_ = response.Serve(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), nil)
		return
	}

	if user == nil {
		uc.lgr.Errorln(fn, tid, response.UserNotFound)
		_ = response.Serve(w, http.StatusNotFound, response.InvalidCredential, nil)
		return
	}

	uc.lgr.Printf(fn, tid, "user.Password: %v, user.Salt: %v", user.Password, user.Salt)
	uc.lgr.Printf(fn, tid, "requested password: %v", loginReq.Password)

	if utils.IsSamePassword(user.Password, loginReq.Password, user.Salt) == false {
		uc.lgr.Errorln(fn, tid, response.InvalidCredential)
		_ = response.Serve(w, http.StatusBadRequest, response.InvalidCredential, nil)
		return
	}

	// generate access token

	// return success response
	_ = response.Serve(w, http.StatusOK, response.Successful, nil)
	return
}

func (uc *UsersController) GetByEmail(w http.ResponseWriter, r *http.Request) {
	fn := "GetUserByEmail"
	tid := utils.GetTracingID(r.Context())

	var getByEmailReq *model.GetUserByEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&getByEmailReq); err != nil {
		_ = response.Serve(w, http.StatusBadRequest, utils.RequiredFieldMessage("email"), nil)
		return
	}

	user, err := uc.svc.GetUserByEmail(r.Context(), &getByEmailReq.Email)
	if err != nil {
		uc.lgr.Errorln(fn, tid, err.Error())
		_ = response.Serve(w, err.StatusCode, err.Error(), nil)
		return
	}

	_ = response.Serve(w, http.StatusOK, response.Successful, user)
	return
}
