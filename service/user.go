package service

import (
	"context"
	"go-rest-api/api/response"
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

func (s *Service) GetByEmail(ctx context.Context, email *string) (*model.UserResponse, *response.Error) {
	fn := "GetByEmail"
	tid := utils.GetTracingID(ctx)
	s.log.Println(fn, tid, "fetching user with email")

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		s.log.Errorln(fn, tid, err.Error())
		return nil, response.NewError(http.StatusNotFound, err.Error())
	}

	s.log.Printf(fn, tid, "user found with email %v", *email)

	return &model.UserResponse{
		UserName: user.UserName,
		Email:    user.Email,
	}, nil
}
