package service

import (
	"fmt"
	"github.com/refandas/scalable-web-service/mygram/core"
	"github.com/refandas/scalable-web-service/mygram/repository"
)

type PhotoService struct {
	photoRepo *repository.PhotoRepo
}

func NewPhotoService(photoRepo *repository.PhotoRepo) *PhotoService {
	return &PhotoService{
		photoRepo: photoRepo,
	}
}

func (s *PhotoService) CreatePhoto(request *core.PhotoCreateRequest) (*core.PhotoResponse, error) {
	user, err := s.photoRepo.CreatePhoto(request)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *PhotoService) FindPhotos() ([]*core.PhotoResponseWithUser, error) {
	photos, err := s.photoRepo.FindPhotos()
	fmt.Println(photos)
	if err != nil {
		return nil, err
	}

	var photoResponses []*core.PhotoResponseWithUser
	for _, photo := range photos {
		photoResponse := core.PhotoResponseWithUser{
			Id:       photo.Id,
			Title:    photo.Title,
			Caption:  *photo.Caption,
			PhotoUrl: photo.PhotoUrl,
			UserId:   photo.UserId,
			User: core.UserResponseSimplified{
				Id:       photo.User.Id,
				Email:    photo.User.Email,
				Username: photo.User.Username,
			},
		}
		photoResponses = append(photoResponses, &photoResponse)
	}

	return photoResponses, err
}

func (s *PhotoService) FindPhotoById(id int) (*core.PhotoResponseWithUser, error) {
	photo, err := s.photoRepo.FindPhotoById(id)
	if err != nil {
		return nil, err
	}

	photoResponse := core.PhotoResponseWithUser{
		Id:       photo.Id,
		Title:    photo.Title,
		Caption:  *photo.Caption,
		PhotoUrl: photo.PhotoUrl,
		UserId:   photo.UserId,
		User: core.UserResponseSimplified{
			Id:       photo.User.Id,
			Email:    photo.User.Email,
			Username: photo.User.Username,
		},
	}

	return &photoResponse, nil
}

func (s *PhotoService) FindPhotoByUserId(userId int) ([]*core.PhotoResponseWithUser, error) {
	photos, err := s.photoRepo.FindPhotosByUserId(userId)
	if err != nil {
		return nil, err
	}

	var photoResponses []*core.PhotoResponseWithUser
	for _, photo := range photos {
		photoResponse := core.PhotoResponseWithUser{
			Id:       photo.Id,
			Title:    photo.Title,
			Caption:  *photo.Caption,
			PhotoUrl: photo.PhotoUrl,
			UserId:   photo.UserId,
			User: core.UserResponseSimplified{
				Id:       photo.User.Id,
				Email:    photo.User.Email,
				Username: photo.User.Username,
			},
		}
		photoResponses = append(photoResponses, &photoResponse)
	}

	return photoResponses, nil
}

func (s *PhotoService) UpdatePhoto(request *core.PhotoUpdateRequest, id int) (*core.PhotoResponse, error) {
	user, err := s.photoRepo.UpdatePhoto(request, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *PhotoService) DeletePhoto(id int) error {
	_, err := s.photoRepo.FindPhotoById(id)
	if err != nil {
		return err
	}

	if err = s.photoRepo.DeletePhoto(id); err != nil {
		return err
	}

	return nil
}
