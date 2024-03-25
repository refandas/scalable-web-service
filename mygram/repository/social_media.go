package repository

import (
	"github.com/refandas/scalable-web-service/mygram/core"
	"github.com/refandas/scalable-web-service/mygram/database"
)

type SocialMediaRepo struct {
	postgres *database.Postgres
}

func NewSocialMediaRepo(db *database.Postgres) *SocialMediaRepo {
	return &SocialMediaRepo{
		postgres: db,
	}
}

func (r *SocialMediaRepo) CreateSocialMedia(request *core.SocialMediaRequest) (*core.SocialMediaResponse, error) {
	var socialMediaResponse core.SocialMediaResponse
	socialMedia := core.SocialMedia{
		Name:           request.Name,
		SocialMediaUrl: request.SocialMediaUrl,
		UserId:         request.UserId,
	}

	if db := r.postgres.DB.Create(&socialMedia).Scan(&socialMediaResponse); db.Error != nil {
		return nil, db.Error
	}

	return &socialMediaResponse, nil
}

func (r *SocialMediaRepo) FindSocialMediaById(id int) (*core.SocialMedia, error) {
	var socialMedia *core.SocialMedia
	if err := r.postgres.DB.Preload("User").First(&socialMedia, id).Error; err != nil {
		return nil, err
	}

	return socialMedia, nil
}

func (r *SocialMediaRepo) FindSocialMediasByUserId(userId int) ([]*core.SocialMedia, error) {
	var socialMedias []*core.SocialMedia
	if err := r.postgres.DB.Preload("User").Where("user_id = ?", userId).Find(&socialMedias).Error; err != nil {
		return nil, err
	}

	return socialMedias, nil
}

func (r *SocialMediaRepo) UpdateSocialMedia(comment *core.SocialMediaRequest, id int) (*core.SocialMediaResponse, error) {
	var socialMediaResponse core.SocialMediaResponse
	db := r.postgres.DB.Model(&core.SocialMedia{}).Where("id = ?", id).Updates(&comment).First(&socialMediaResponse)
	if db.Error != nil {
		return nil, db.Error
	}

	return &socialMediaResponse, nil
}

func (r *SocialMediaRepo) DeleteSocialMedia(id int) error {
	if db := r.postgres.DB.Delete(&core.SocialMedia{}, id); db.Error != nil {
		return db.Error
	}

	return nil
}
