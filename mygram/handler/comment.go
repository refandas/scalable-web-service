package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/refandas/scalable-web-service/mygram/core"
	"github.com/refandas/scalable-web-service/mygram/helper"
	"github.com/refandas/scalable-web-service/mygram/service"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

type CommentHandler struct {
	commentService *service.CommentService
}

func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

func (h *CommentHandler) CreateComment(ctx *gin.Context) {
	var request core.CommentCreateRequest

	if err := ctx.BindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
		})
		return
	}

	validator := helper.NewValidator()

	if err := validator.Validate(request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request format",
		})
		return
	}

	user := ctx.MustGet("userData").(jwt.MapClaims)
	request.UserId = int64(user["id"].(float64))

	res, err := h.commentService.CreateComment(&request)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "Photo not found",
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create comment",
		})
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

func (h *CommentHandler) FindComments(ctx *gin.Context) {
	comments, err := h.commentService.FindComments()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch comments",
		})
	}

	if comments == nil {
		comments = []*core.CommentResponseWithUserAndPhoto{}
	}

	ctx.JSON(http.StatusOK, comments)
}

func (h *CommentHandler) FindCommentById(ctx *gin.Context) {
	idStr := ctx.Param("commentId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid comment id",
		})
	}

	photo, err := h.commentService.FindCommentById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "Comment not found",
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch comment",
		})
		return
	}

	ctx.JSON(http.StatusOK, photo)
}

func (h *CommentHandler) UpdateComment(ctx *gin.Context) {
	idStr := ctx.Param("commentId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid comment id",
		})
	}

	var request *core.CommentUpdateRequest
	if err := ctx.BindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
		})
		return
	}

	validator := helper.NewValidator()

	if err := validator.Validate(request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request format",
		})
		return
	}

	res, err := h.commentService.UpdateComment(request, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "Comment not found",
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update comment",
		})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (h *CommentHandler) DeleteComment(ctx *gin.Context) {
	idStr := ctx.Param("commentId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid comment id",
		})
	}

	err = h.commentService.DeleteComment(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "Comment not found",
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Cannot delete comment",
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}
