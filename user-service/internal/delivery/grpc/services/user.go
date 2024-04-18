package services

import (
	"context"
	"time"
	userproto "user-service/genproto/user_service"
	"user-service/internal/entity"
	"user-service/internal/usecase"
	"user-service/internal/usecase/event"

	"go.uber.org/zap"
)

type userRPC struct {
	logger         *zap.Logger
	userUsecase    usecase.User
	brokerProducer event.BrokerProducer
}

func NewRPC(logger *zap.Logger, userUsecase usecase.User, brokerProducer event.BrokerProducer) userproto.UserServiceServer {
	return &userRPC{
		logger:         logger,
		userUsecase:    userUsecase,
		brokerProducer: brokerProducer,
	}
}

func (s userRPC) Create(ctx context.Context, in *userproto.User) (*userproto.GetUserRequest, error) {
	guid, err := s.userUsecase.Create(ctx, &entity.User{
		GUID:      in.Id,
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Username:  in.Username,
		Email:     in.Email,
		Password:  in.Password,
		Bio:       in.Bio,
		Website:   in.Website,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return &userproto.GetUserRequest{
		Id: guid,
	}, nil
}

func (s userRPC) Update(ctx context.Context, in *userproto.User) (*userproto.User, error) {
	err := s.userUsecase.Update(ctx, &entity.User{
		GUID:      in.Id,
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Username:  in.Username,
		Email:     in.Email,
		Password:  in.Password,
		Bio:       in.Password,
		Website:   in.Website,
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	return &userproto.User{
		Id:        in.Id,
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Username:  in.Username,
		Email:     in.Email,
		Password:  in.Password,
		Bio:       in.Bio,
		Website:   in.Website,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
	}, nil
}

func (s userRPC) Delete(ctx context.Context, in *userproto.GetUserRequest) (*userproto.DeletedUser, error) {
	if err := s.userUsecase.Delete(ctx, in.Id); err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	return &userproto.DeletedUser{}, nil
}

func (s userRPC) Get(ctx context.Context, in *userproto.GetUserRequest) (*userproto.User, error) {
	user, err := s.userUsecase.Get(ctx, map[string]string{
		"id": in.Id,
	})

	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return &userproto.User{
		Id:        user.GUID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		Bio:       user.Bio,
		Website:   user.Website,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s userRPC) GetAll(ctx context.Context, in *userproto.GetAllUserRequest) (*userproto.GetAllUserResponse, error) {
	offset := in.Limit * (in.Page - 1)
	users, err := s.userUsecase.List(ctx, uint64(in.Limit), uint64(offset), map[string]string{})
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	var response userproto.GetAllUserResponse
	for _, u := range users {

		temp := &userproto.User{
			Id:        u.GUID,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Username:  u.Username,
			Email:     u.Email,
			Password:  u.Password,
			Bio:       u.Bio,
			Website:   u.Website,
			CreatedAt: u.CreatedAt.Format(time.RFC3339),
			UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
		}

		response.AllUsers = append(response.AllUsers, temp)
	}

	return &response, nil
}
