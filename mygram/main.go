package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/refandas/scalable-web-service/mygram/database"
	"github.com/refandas/scalable-web-service/mygram/handler"
	"github.com/refandas/scalable-web-service/mygram/middleware"
	"github.com/refandas/scalable-web-service/mygram/repository"
	"github.com/refandas/scalable-web-service/mygram/service"
	"log"
	"net/http"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	postgres := database.NewPostgres()
	if postgres.Err != nil {
		panic(postgres.Err)
	}

	r := gin.Default()

	r.GET("", rootHandler)

	userRepo := repository.NewUserRepo(postgres)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	users := r.Group("/users")

	users.POST("/register", userHandler.CreateUser)
	users.POST("/login", userHandler.Login)
	users.Use(middleware.Authentication())
	{
		users.PUT("/", userHandler.UpdateUser)
		users.DELETE("/", userHandler.DeleteUser)
	}

	photoRepo := repository.NewPhotoRepo(postgres)
	photoService := service.NewPhotoService(photoRepo)
	photoHandler := handler.NewPhotoHandler(photoService)

	photos := r.Group("/photos")
	photos.Use(middleware.Authentication())
	{
		photos.GET("/", photoHandler.FindPhotos)
		photos.GET("/:photoId", photoHandler.FindPhotoById)
		photos.POST("/", photoHandler.CreatePhoto)
		photos.PUT("/:photoId", photoHandler.UpdatePhoto)
		photos.DELETE("/:photoId", photoHandler.DeletePhoto)
	}

	commentRepo := repository.NewCommentRepo(postgres)
	commentService := service.NewCommentService(commentRepo)
	commentHandler := handler.NewCommentHandler(commentService)

	comments := r.Group("/comments")
	comments.Use(middleware.Authentication())
	{
		comments.POST("/", commentHandler.CreateComment)
		comments.GET("/", commentHandler.FindComments)
		comments.GET("/:commentId", commentHandler.FindCommentById)
		comments.PUT("/:commentId", commentHandler.UpdateComment)
		comments.DELETE("/:commentId", commentHandler.DeleteComment)
	}

	socialMediaRepo := repository.NewSocialMediaRepo(postgres)
	socialMediaService := service.NewSocialMedia(socialMediaRepo)
	socialMediaHandler := handler.NewSocialMediaHandler(socialMediaService)

	socialMedias := r.Group("/socialmedias")
	socialMedias.Use(middleware.Authentication())
	{
		socialMedias.POST("/", socialMediaHandler.CreateSocialMedia)
		socialMedias.GET("/", socialMediaHandler.FindSocialMedias)
		socialMedias.GET("/:socialMediaId", socialMediaHandler.FindSocialMediaById)
		socialMedias.PUT("/:socialMediaId", socialMediaHandler.UpdateSocialMedia)
		socialMedias.DELETE("/:socialMediaId", socialMediaHandler.DeleteSocialMedia)
	}

	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func rootHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "Service is up and running",
	})
}
