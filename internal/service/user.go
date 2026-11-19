package service

import (
	"ca/internal/model"
	"ca/internal/repository/postgres"
	"ca/internal/transport/dto"
	"context"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository postgres.UserRepository
}

func NewUserService(userRepository postgres.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) Login(ctx context.Context, data *dto.LoginRequest) (model.User, error) {
	user, err := s.userRepository.GetByUsername(ctx, data.Username)
	if err != nil {
		return model.User{}, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		return model.User{}, err
	}

	return *user, nil
}

func (s *UserService) GetUserById(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return s.userRepository.GetById(ctx, id)
}

func (s *UserService) Register(ctx context.Context, data *dto.UserCreateRequest) (id uuid.UUID, err error) {
	return s.userRepository.Create(ctx, data)
}

func (s *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.userRepository.Delete(ctx, id)
}
