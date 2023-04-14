package service

import (
	"context"
	"go-rest-api/api/response"
	"go-rest-api/infra"
	"go-rest-api/model"
	"go-rest-api/utils"
	"net/http"
)

func (s *Service) CreateUser(ctx context.Context, user *model.User) *response.Error {
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

func (s *Service) GetUserByEmail(ctx context.Context, email *string) (*model.GetUserByEmailResponse, *response.Error) {
	fn := "GetUserByEmail"
	tid := utils.GetTracingID(ctx)
	s.log.Println(fn, tid, "fetching user with email")

	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		s.log.Errorln(fn, tid, err.Error())
		return nil, response.NewError(http.StatusNotFound, err.Error())
	}

	s.log.Printf(fn, tid, "user found with email %v", *email)

	return &model.GetUserByEmailResponse{
		UserName: user.UserName,
		Email:    user.Email,
	}, nil
}

func (s *Service) GetAuthUserByEmail(ctx context.Context, email *string) (*model.AuthUserPrivateData, error) {
	fn := "GetAuthUserByEmail"
	tid := utils.GetTracingID(ctx)

	s.log.Println(fn, tid, "fetching auth user by email")

	user, err := s.userRepo.GetAuthUserByEmail(ctx, email)
	if err != nil {
		s.log.Errorln(fn, tid, err.Error())
		return nil, infra.ErrNotFound
	}

	s.log.Printf(fn, tid, "user found with email %v", *email)

	return user, nil
}
