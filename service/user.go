package service

import (
	"context"
	"go-rest-api/api/response"
	"go-rest-api/infra"
	"go-rest-api/logger"
	"go-rest-api/model"
	"go-rest-api/repo"
	"go-rest-api/utils"
	"net/http"
)

type UserService interface {
	CreateUser(ctx context.Context, user *model.User) *response.Error
	GetUserByEmail(ctx context.Context, email *string) (*model.GetUserByEmailResponse, *response.Error)
	GetAuthUserByEmail(ctx context.Context, email *string) (*model.AuthUserData, error)
	StoreToken(ctx context.Context, email string, token string) error
	GetToken(ctx context.Context, key string) (string, error)
	RemoveToken(ctx context.Context, key ...string) error
}

type userService struct {
	log      logger.StructLogger
	userRepo repo.UserRepo
}

func NewUserService(lgr logger.StructLogger, userRepo repo.UserRepo) UserService {
	return &userService{
		log:      lgr,
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *model.User) *response.Error {
	fn := "CreateUser"
	tid := utils.GetTracingID(ctx)
	s.log.Println(fn, tid, "creating user")

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.log.Errorln(fn, tid, err.Error())
		return response.NewError(http.StatusInternalServerError, err.Error())
	}

	s.log.Println(fn, tid, "user created successfully")
	return nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email *string) (*model.GetUserByEmailResponse, *response.Error) {
	fn := "GetUserByEmail"
	tid := utils.GetTracingID(ctx)
	s.log.Println(fn, tid, "fetching user with email")

	user, err := s.userRepo.GetUserByEmail(ctx, *email)
	if err != nil {
		s.log.Errorln(fn, tid, err.Error())
		return nil, response.NewError(http.StatusNotFound, err.Error())
	}

	s.log.Printf(fn, tid, "user found with email %v", *email)

	return user, nil
}

func (s *userService) GetAuthUserByEmail(ctx context.Context, email *string) (*model.AuthUserData, error) {
	fn := "GetAuthUserByEmail"
	tid := utils.GetTracingID(ctx)

	s.log.Println(fn, tid, "fetching auth user by email")

	user, err := s.userRepo.GetAuthUserByEmail(ctx, *email)
	if err != nil {
		s.log.Errorln(fn, tid, err.Error())
		return nil, infra.ErrNotFound
	}

	s.log.Printf(fn, tid, "user found with email %v", *email)

	return user, nil
}

func (s *userService) StoreToken(ctx context.Context, email string, token string) error {
	fn := "StoreToken"
	tid := utils.GetTracingID(ctx)

	s.log.Println(fn, tid, "storing token")

	if err := s.userRepo.StoreToken(ctx, email, token); err != nil {
		s.log.Errorln(fn, tid, err.Error())
		return err
	}
	return nil
}

func (s *userService) GetToken(ctx context.Context, key string) (string, error) {
	fn := "GetToken"
	tid := utils.GetTracingID(ctx)

	s.log.Println(fn, tid, "fetching token")

	value, err := s.userRepo.GetToken(ctx, key)

	if err != nil {
		s.log.Errorln(fn, tid, err.Error())
		return "", err
	}
	return value, err
}

func (s *userService) RemoveToken(ctx context.Context, key ...string) error {
	fn := "RemoveToken"
	tid := utils.GetTracingID(ctx)

	s.log.Println(fn, tid, "removing token")

	if err := s.userRepo.RemoveToken(ctx, key...); err != nil {
		s.log.Errorln(fn, tid, err.Error())
		return err
	}

	return nil
}
