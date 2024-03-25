package service

import (
	"github.com/refandas/scalable-web-service/mygram/core"
	"github.com/refandas/scalable-web-service/mygram/repository"
)

type CommentService struct {
	commentRepo *repository.CommentRepo
}

func NewCommentService(commentRepo *repository.CommentRepo) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
	}
}

func (s *CommentService) CreateComment(request *core.CommentCreateRequest) (*core.CommentResponse, error) {
	user, err := s.commentRepo.CreateComment(request)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *CommentService) FindComments() ([]*core.CommentResponseWithUserAndPhoto, error) {
	comments, err := s.commentRepo.FindComments()
	if err != nil {
		return nil, err
	}

	var commentResponses []*core.CommentResponseWithUserAndPhoto
	for _, comment := range comments {
		commentResponse := core.CommentResponseWithUserAndPhoto{
			Id:      comment.Id,
			Message: comment.Message,
			PhotoId: comment.PhotoId,
			UserId:  comment.UserId,
			User: core.UserResponseSimplified{
				Id:       comment.User.Id,
				Email:    comment.User.Email,
				Username: comment.User.Username,
			},
			Photo: core.PhotoResponse{
				Id:       comment.Photo.Id,
				Title:    comment.Photo.Title,
				Caption:  *comment.Photo.Caption,
				PhotoUrl: comment.Photo.PhotoUrl,
				UserId:   comment.Photo.UserId,
			},
		}
		commentResponses = append(commentResponses, &commentResponse)
	}

	return commentResponses, err
}

func (s *CommentService) FindCommentById(id int) (*core.CommentResponseWithUserAndPhoto, error) {
	comment, err := s.commentRepo.FindCommentById(id)
	if err != nil {
		return nil, err
	}

	commentResponse := core.CommentResponseWithUserAndPhoto{
		Id:      comment.Id,
		Message: comment.Message,
		PhotoId: comment.PhotoId,
		UserId:  comment.UserId,
		User: core.UserResponseSimplified{
			Id:       comment.User.Id,
			Email:    comment.User.Email,
			Username: comment.User.Username,
		},
		Photo: core.PhotoResponse{
			Id:       comment.Photo.Id,
			Title:    comment.Photo.Title,
			Caption:  *comment.Photo.Caption,
			PhotoUrl: comment.Photo.PhotoUrl,
			UserId:   comment.Photo.UserId,
		},
	}

	return &commentResponse, nil
}

func (s *CommentService) FindCommentByPhotoId(photoId int) ([]*core.CommentResponse, error) {
	comments, err := s.commentRepo.FindCommentsByPhotoId(photoId)
	if err != nil {
		return nil, err
	}

	var commentResponses []*core.CommentResponse
	for _, comment := range comments {
		commentResponse := core.CommentResponse{
			Id:      comment.Id,
			Message: comment.Message,
			PhotoId: comment.PhotoId,
			UserId:  comment.UserId,
		}
		commentResponses = append(commentResponses, &commentResponse)
	}

	return commentResponses, err
}

func (s *CommentService) UpdateComment(request *core.CommentUpdateRequest, id int) (*core.CommentResponse, error) {
	user, err := s.commentRepo.UpdateComments(request, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *CommentService) DeleteComment(id int) error {
	_, err := s.commentRepo.FindCommentById(id)
	if err != nil {
		return err
	}

	if err = s.commentRepo.DeleteComment(id); err != nil {
		return err
	}

	return nil
}
