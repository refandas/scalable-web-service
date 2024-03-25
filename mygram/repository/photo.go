package repository

import (
	"github.com/refandas/scalable-web-service/mygram/core"
	"github.com/refandas/scalable-web-service/mygram/database"
)

type PhotoRepo struct {
	postgres *database.Postgres
}

func NewPhotoRepo(db *database.Postgres) *PhotoRepo {
	return &PhotoRepo{
		postgres: db,
	}
}

func (r *PhotoRepo) CreatePhoto(request *core.PhotoCreateRequest) (*core.PhotoResponse, error) {
	var photoResponse core.PhotoResponse
	photo := core.Photo{
		Title:    request.Title,
		Caption:  &request.Caption,
		PhotoUrl: request.PhotoUrl,
		UserId:   request.UserId,
	}

	if db := r.postgres.DB.Create(&photo).Scan(&photoResponse); db.Error != nil {
		return nil, db.Error
	}

	return &photoResponse, nil
}

func (r *PhotoRepo) FindPhotos() ([]*core.Photo, error) {
	var photos []*core.Photo
	if err := r.postgres.DB.Preload("User").Find(&photos).Error; err != nil {
		return nil, err
	}

	return photos, nil
}

func (r *PhotoRepo) FindPhotosByUserId(userId int) ([]*core.Photo, error) {
	var photos []*core.Photo
	if err := r.postgres.DB.Preload("User").Where("user_id = ?", userId).Find(&photos).Error; err != nil {
		return nil, err
	}

	return photos, nil
}

func (r *PhotoRepo) FindPhotoById(id int) (*core.Photo, error) {
	var photo core.Photo
	if err := r.postgres.DB.Preload("User").First(&photo, id).Error; err != nil {
		return nil, err
	}

	return &photo, nil
}

func (r *PhotoRepo) UpdatePhoto(photo *core.PhotoUpdateRequest, id int) (*core.PhotoResponse, error) {
	var photoResponse core.PhotoResponse
	db := r.postgres.DB.Model(&core.Photo{}).Where("id = ?", id).Updates(&photo).First(&photoResponse)
	if db.Error != nil {
		return nil, db.Error
	}

	return &photoResponse, nil
}

func (r *PhotoRepo) DeletePhoto(id int) error {
	if db := r.postgres.DB.Delete(&core.Photo{}, id); db.Error != nil {
		return db.Error
	}

	return nil
}
