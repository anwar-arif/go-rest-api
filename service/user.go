package service

import (
	"context"
	"go-rest-api/api/response"
	"go-rest-api/model"
	"go-rest-api/utils"
	"net/http"
)

func (s *Service) CreateUser(ctx context.Context, user *model.User) *response.Error {
	tid := utils.GetTracingID(ctx)
	s.log.Println("CreateUser", tid, "creating user")

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.log.Errorln("CreateUser", tid, err.Error())
		return response.NewError(http.StatusInternalServerError, err.Error())
	}

	s.log.Println("CreateUser", tid, "user created successfully")
	return nil
}

func (s *Service) GetByEmail(ctx context.Context, email *string) (*model.User, *response.Error) {
	tid := utils.GetTracingID(ctx)
	s.log.Println("GetByEmail", tid, "fetching user with email")

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		s.log.Errorln("GetByEmail", tid, err.Error())
		return nil, response.NewError(http.StatusNotFound, err.Error())
	}

	s.log.Println("GetByEmail", tid, "user found")
	return user, nil
}
