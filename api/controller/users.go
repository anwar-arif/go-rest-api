package controller

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go-rest-api/api/response"
	"go-rest-api/config"
	"go-rest-api/infra"
	infra2 "go-rest-api/infra/db"
	"go-rest-api/logger"
	"go-rest-api/model"
	"go-rest-api/repo"
	"go-rest-api/service"
	"go-rest-api/utils"
	"net/http"
)

type UsersController interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	LogIn(w http.ResponseWriter, r *http.Request)
	LogOut(w http.ResponseWriter, r *http.Request)
	GetByEmail(w http.ResponseWriter, r *http.Request)
}

type usersController struct {
	svc service.UserService
	lgr logger.StructLogger
}

func NewUsersController(cfgDBTable *config.Table, db *infra2.DB, lgr logger.StructLogger) UsersController {
	userRepo := repo.NewUser(cfgDBTable.UserCollectionName, db)
	userService := service.NewUserService(lgr, userRepo)

	return &usersController{
		svc: userService,
		lgr: lgr,
	}
}

func (uc *usersController) CreateUser(w http.ResponseWriter, r *http.Request) {
	fn := "CreateUser"
	tid := utils.GetTracingID(r.Context())

	var signUpReq *model.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&signUpReq); err != nil {
		uc.lgr.Errorf(fn, tid, "error while parsing user payload: %v\n", err.Error())
		_ = response.Serve(w, http.StatusBadRequest, utils.RequiredFieldMessage(), nil)
		return
	}

	// user exists with same email
	existingUser, err := uc.svc.GetUserByEmail(r.Context(), &signUpReq.Email)
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

	u := &model.User{
		UserName: signUpReq.UserName,
		Email:    signUpReq.Email,
	}

	u.Salt = *salt
	u.Password = utils.HashPassword(signUpReq.Password, u.Salt)
	u.UserID = uuid.New().String()

	uc.lgr.Printf(fn, tid, "user id: ", u.UserID)

	if err := uc.svc.CreateUser(r.Context(), u); err != nil {
		uc.lgr.Errorln(fn, tid, err.Error())
		_ = response.Serve(w, err.StatusCode, err.Error(), nil)
		return
	}

	_ = response.Serve(w, http.StatusCreated, http.StatusText(http.StatusCreated), &model.SignUpResponse{
		UserName: u.UserName,
		Email:    u.Email,
	})
	return
}

func (uc *usersController) LogIn(w http.ResponseWriter, r *http.Request) {
	fn := "LogIn"
	tid := utils.GetTracingID(r.Context())

	var loginReq *model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		_ = response.Serve(w, http.StatusBadRequest, response.CannotProcessRequest, nil)
		return
	}

	au, err := uc.svc.GetAuthUserByEmail(r.Context(), &loginReq.Email)
	if (err != nil) && (err.Error() != infra.ErrNotFound.Error()) {
		uc.lgr.Errorln(fn, tid, http.StatusText(http.StatusInternalServerError))
		_ = response.Serve(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), nil)
		return
	}

	if au == nil {
		uc.lgr.Errorln(fn, tid, response.UserNotFound)
		_ = response.Serve(w, http.StatusNotFound, response.UserNotFound, nil)
		return
	}

	if utils.IsSamePassword(au.Password, loginReq.Password, au.Salt) == false {
		uc.lgr.Errorln(fn, tid, response.InvalidCredential)
		_ = response.Serve(w, http.StatusBadRequest, response.InvalidCredential, nil)
		return
	}

	// generate access token
	viper.AutomaticEnv()
	secretKey := viper.GetString("app.api_secret_key")
	expiresIn := viper.GetInt64("app.access_token_expires_in")
	uc.lgr.Printf(fn, tid, "expiresIn: %v", expiresIn)

	token, err := utils.GenerateToken(au.ToUser(), secretKey)

	// store token
	if err := uc.svc.StoreToken(r.Context(), au.Email, token); err != nil {
		uc.lgr.Println(fn, tid, utils.TokenCantStore)
		_ = response.Serve(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), nil)
		return
	}
	// return success response
	_ = response.Serve(w, http.StatusOK, response.Successful, model.LoginResponse{
		Email:       au.Email,
		AccessToken: token,
	})
	return
}

func (uc *usersController) LogOut(w http.ResponseWriter, r *http.Request) {
	fn := "LogOut"
	tid := utils.GetTracingID(r.Context())

	var logOutReq *model.LogOutRequest
	if err := json.NewDecoder(r.Body).Decode(&logOutReq); err != nil {
		uc.lgr.Errorln(fn, tid, err.Error())
		_ = response.Serve(w, http.StatusBadRequest, utils.RequiredFieldMessage("email"), nil)
		return
	}

	_, cerr := utils.ClaimsFromRequest(r)
	if cerr != nil {
		uc.lgr.Errorln(fn, tid, cerr.Error())
		_ = response.Serve(w, http.StatusBadRequest, utils.CantProcessRequest, nil)
		return
	}

	if err := uc.svc.RemoveToken(r.Context(), logOutReq.Email); err != nil {
		uc.lgr.Errorln(fn, tid, err.Error())
		_ = response.Serve(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), nil)
		return
	}

	_ = response.Serve(w, http.StatusOK, utils.LogoutSuccessful, nil)
	return
}

func (uc *usersController) GetByEmail(w http.ResponseWriter, r *http.Request) {
	fn := "GetUserByEmail"
	tid := utils.GetTracingID(r.Context())

	var getByEmailReq *model.GetUserByEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&getByEmailReq); err != nil {
		_ = response.Serve(w, http.StatusBadRequest, utils.RequiredFieldMessage("email"), nil)
		return
	}
	claims, claimsErr := utils.ClaimsFromRequest(r)
	if (claimsErr != nil) || (claims.Email != getByEmailReq.Email) {
		if claims != nil {
			uc.lgr.Printf(fn, tid, claims.Email)
		}
		uc.lgr.Errorln(fn, tid, fmt.Sprintf("%s: %s", utils.UnAuthorizedAccess, claimsErr.Error()))
		_ = response.Serve(w, http.StatusUnauthorized, utils.UnAuthorizedAccess, nil)
		return
	}

	tokenInReq := utils.GetAuthTokenFromHeader(r)

	// returns err when key doesn't exist
	storedToken, tokenErr := uc.svc.GetToken(r.Context(), claims.Email)

	if (tokenErr != nil) || (storedToken != tokenInReq) {
		uc.lgr.Errorln(fn, tid, fmt.Sprintf("%s: %s", utils.TokenNotFound, tokenErr.Error()))
		_ = response.Serve(w, http.StatusUnauthorized, utils.UnAuthorizedAccess, nil)
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
