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
)

type PhotoHandler struct {
	photoService *service.PhotoService
}

func NewPhotoHandler(photoService *service.PhotoService) *PhotoHandler {
	return &PhotoHandler{
		photoService: photoService,
	}
}

func (h *PhotoHandler) CreatePhoto(ctx *gin.Context) {
	var request core.PhotoCreateRequest

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

	res, err := h.photoService.CreatePhoto(&request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create photo",
		})
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

func (h *PhotoHandler) FindPhotos(ctx *gin.Context) {
	photos, err := h.photoService.FindPhotos()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch photos",
		})
	}

	if photos == nil {
		photos = []*core.PhotoResponseWithUser{}
	}

	ctx.JSON(http.StatusOK, photos)
}

func (h *PhotoHandler) FindPhotoById(ctx *gin.Context) {
	idStr := ctx.Param("photoId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid photo id",
		})
	}

	photo, err := h.photoService.FindPhotoById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "Photo not found",
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch photo",
		})
		return
	}

	ctx.JSON(http.StatusOK, photo)
}

func (h *PhotoHandler) UpdatePhoto(ctx *gin.Context) {
	idStr := ctx.Param("photoId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid photo id",
		})
	}

	var request *core.PhotoUpdateRequest
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

	res, err := h.photoService.UpdatePhoto(request, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "Photo not found",
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update photo",
		})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (h *PhotoHandler) DeletePhoto(ctx *gin.Context) {
	idStr := ctx.Param("photoId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid photo id",
		})
	}

	err = h.photoService.DeletePhoto(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "Photo not found",
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Cannot delete photo",
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}
