package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang-jwt/jwt/v5"
	"github.com/refandas/scalable-web-service/mygram/core"
	"github.com/refandas/scalable-web-service/mygram/helper"
	"github.com/refandas/scalable-web-service/mygram/service"
	"net/http"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(ctx *gin.Context) {
	var request core.UserCreateRequest

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

	hashedPassword := helper.HashPassword(request.Password)
	request.Password = hashedPassword

	res, err := h.userService.CreateUser(&request)
	if err != nil {
		if err.Error() == "duplicate email" {
			ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"message": "This email is already in use",
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create user",
		})
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	var request *core.UserUpdateRequest
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
	res, err := h.userService.UpdateUser(request, int(user["id"].(float64)))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update user",
		})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	user := ctx.MustGet("userData").(jwt.MapClaims)

	err := h.userService.DeleteUser(int(user["id"].(float64)))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Cannot delete user",
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *UserHandler) Login(ctx *gin.Context) {
	request := &core.LoginRequest{}

	err := ctx.ShouldBindBodyWith(request, binding.JSON)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
		})
		return
	}

	validator := helper.NewValidator()

	err = validator.Validate(request)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request format",
		})
		return
	}

	user, err := h.userService.FindUserByEmail(request.Email)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "Wrong email/password",
		})
		return
	}

	result := helper.ValidateHashPassword(request.Password, user.Password)
	if !result {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "Wrong email/password",
		})
		return
	}

	token, err := helper.GenerateToken(user.Id, user.Email)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create token",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
