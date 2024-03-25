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

type SocialMediaHandler struct {
	socialMediaService *service.SocialMediaService
}

func NewSocialMediaHandler(socialMediaService *service.SocialMediaService) *SocialMediaHandler {
	return &SocialMediaHandler{
		socialMediaService: socialMediaService,
	}
}

func (h *SocialMediaHandler) CreateSocialMedia(ctx *gin.Context) {
	var request core.SocialMediaRequest

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

	res, err := h.socialMediaService.CreateSocialMedia(&request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create comment",
		})
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

func (h *SocialMediaHandler) FindSocialMedias(ctx *gin.Context) {
	user := ctx.MustGet("userData").(jwt.MapClaims)

	socialMedias, err := h.socialMediaService.FindSocialMediasByUserId(int(user["id"].(float64)))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch social medias",
		})
	}

	if socialMedias == nil {
		socialMedias = []*core.SocialMediaResponseWithUser{}
	}

	ctx.JSON(http.StatusOK, socialMedias)
}

func (h *SocialMediaHandler) FindSocialMediaById(ctx *gin.Context) {
	idStr := ctx.Param("socialMediaId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid social media id",
		})
	}

	socialMedia, err := h.socialMediaService.FindSocialMediaById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "Social media not found",
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch social media",
		})
		return
	}

	ctx.JSON(http.StatusOK, socialMedia)
}

func (h *SocialMediaHandler) UpdateSocialMedia(ctx *gin.Context) {
	idStr := ctx.Param("socialMediaId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid social media id",
		})
	}

	var request *core.SocialMediaRequest
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

	res, err := h.socialMediaService.UpdateSocialMedia(request, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "Social media not found",
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update social media",
		})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (h *SocialMediaHandler) DeleteSocialMedia(ctx *gin.Context) {
	idStr := ctx.Param("socialMediaId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid social media id",
		})
	}

	err = h.socialMediaService.DeleteSocialMedia(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": "Social media not found",
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Cannot delete social media",
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}
