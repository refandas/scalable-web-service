package repository

import (
	"errors"
	"fmt"
	"github.com/refandas/scalable-web-service/mygram/core"
	"github.com/refandas/scalable-web-service/mygram/database"
	"gorm.io/gorm"
)

type CommentRepo struct {
	postgres *database.Postgres
}

func NewCommentRepo(db *database.Postgres) *CommentRepo {
	return &CommentRepo{
		postgres: db,
	}
}

func (r *CommentRepo) CreateComment(request *core.CommentCreateRequest) (*core.CommentResponse, error) {
	var photo core.Photo
	if err := r.postgres.DB.First(&photo, request.PhotoId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("photo not found")
		}
		return nil, err
	}

	var commentResponse core.CommentResponse
	comment := core.Comment{
		UserId:  request.UserId,
		PhotoId: request.PhotoId,
		Message: request.Message,
	}

	if db := r.postgres.DB.Create(&comment).Scan(&commentResponse); db.Error != nil {
		return nil, db.Error
	}

	return &commentResponse, nil
}

func (r *CommentRepo) FindComments() ([]*core.Comment, error) {
	var comments []*core.Comment
	if err := r.postgres.DB.Preload("User").Preload("Photo").Find(&comments).Error; err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *CommentRepo) FindCommentById(id int) (*core.Comment, error) {
	var comment *core.Comment
	if err := r.postgres.DB.Preload("User").Preload("Photo").First(&comment, id).Error; err != nil {
		return nil, err
	}

	return comment, nil
}

func (r *CommentRepo) FindCommentsByPhotoId(photoId int) ([]*core.Comment, error) {
	var comments []*core.Comment
	if err := r.postgres.DB.Where("photo_id = ?", photoId).Find(&comments).Error; err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *CommentRepo) UpdateComments(comment *core.CommentUpdateRequest, id int) (*core.CommentResponse, error) {
	var commentResponse core.CommentResponse
	db := r.postgres.DB.Model(&core.Comment{}).Where("id = ?", id).Updates(&comment).First(&commentResponse)
	if db.Error != nil {
		return nil, db.Error
	}

	return &commentResponse, nil
}

func (r *CommentRepo) DeleteComment(id int) error {
	if db := r.postgres.DB.Delete(&core.Comment{}, id); db.Error != nil {
		return db.Error
	}

	return nil
}
