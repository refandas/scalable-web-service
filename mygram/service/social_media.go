package service

import (
	"github.com/refandas/scalable-web-service/mygram/core"
	"github.com/refandas/scalable-web-service/mygram/repository"
)

type SocialMediaService struct {
	socialMediaRepo *repository.SocialMediaRepo
}

func NewSocialMedia(socialMediaRepo *repository.SocialMediaRepo) *SocialMediaService {
	return &SocialMediaService{
		socialMediaRepo: socialMediaRepo,
	}
}

func (s *SocialMediaService) CreateSocialMedia(request *core.SocialMediaRequest) (*core.SocialMediaResponse, error) {
	socialMedia, err := s.socialMediaRepo.CreateSocialMedia(request)
	if err != nil {
		return nil, err
	}

	return socialMedia, nil
}

func (s *SocialMediaService) FindSocialMediasByUserId(userId int) ([]*core.SocialMediaResponseWithUser, error) {
	socialMedias, err := s.socialMediaRepo.FindSocialMediasByUserId(userId)
	if err != nil {
		return nil, err
	}

	var socialMediaResponses []*core.SocialMediaResponseWithUser
	for _, socialMedia := range socialMedias {
		commentResponse := core.SocialMediaResponseWithUser{
			Id:             socialMedia.Id,
			Name:           socialMedia.Name,
			SocialMediaUrl: socialMedia.SocialMediaUrl,
			UserId:         socialMedia.UserId,
			User: core.UserResponseSimplified{
				Id:       socialMedia.User.Id,
				Email:    socialMedia.User.Email,
				Username: socialMedia.User.Username,
			},
		}
		socialMediaResponses = append(socialMediaResponses, &commentResponse)
	}

	return socialMediaResponses, err
}

func (s *SocialMediaService) FindSocialMediaById(id int) (*core.SocialMediaResponseWithUser, error) {
	socialMedia, err := s.socialMediaRepo.FindSocialMediaById(id)
	if err != nil {
		return nil, err
	}

	socialMediaResponse := core.SocialMediaResponseWithUser{
		Id:             socialMedia.Id,
		Name:           socialMedia.Name,
		SocialMediaUrl: socialMedia.SocialMediaUrl,
		UserId:         socialMedia.UserId,
		User: core.UserResponseSimplified{
			Id:       socialMedia.User.Id,
			Email:    socialMedia.User.Email,
			Username: socialMedia.User.Username,
		},
	}

	return &socialMediaResponse, nil
}

func (s *SocialMediaService) UpdateSocialMedia(request *core.SocialMediaRequest, id int) (*core.SocialMediaResponse, error) {
	user, err := s.socialMediaRepo.UpdateSocialMedia(request, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *SocialMediaService) DeleteSocialMedia(id int) error {
	_, err := s.socialMediaRepo.FindSocialMediaById(id)
	if err != nil {
		return err
	}

	if err = s.socialMediaRepo.DeleteSocialMedia(id); err != nil {
		return err
	}

	return nil
}
