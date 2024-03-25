package service

import (
	"github.com/refandas/scalable-web-service/mygram/core"
	"github.com/refandas/scalable-web-service/mygram/repository"
)

type UserService struct {
	userRepo *repository.UserRepo
}

func NewUserService(userRepo *repository.UserRepo) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUser(request *core.UserCreateRequest) (*core.UserResponse, error) {
	user, err := s.userRepo.CreateUser(request)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) FindUserById(id int) (*core.User, error) {
	user, err := s.userRepo.FindUserById(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) FindUserByEmail(email string) (*core.User, error) {
	user, err := s.userRepo.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateUser(request *core.UserUpdateRequest, id int) (*core.UserResponse, error) {
	user, err := s.userRepo.UpdateUser(request, id)
	if err != nil {
		return nil, err
	}

	return user, nil

}

func (s *UserService) DeleteUser(id int) error {
	if err := s.userRepo.DeleteUser(id); err != nil {
		return err
	}

	return nil
}
