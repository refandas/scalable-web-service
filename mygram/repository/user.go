package repository

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/refandas/scalable-web-service/mygram/core"
	"github.com/refandas/scalable-web-service/mygram/database"
)

type UserRepo struct {
	postgres *database.Postgres
}

func NewUserRepo(db *database.Postgres) *UserRepo {
	return &UserRepo{
		postgres: db,
	}
}

func (r *UserRepo) CreateUser(request *core.UserCreateRequest) (*core.UserResponse, error) {
	var userResponse core.UserResponse
	user := core.User{
		Email:           request.Email,
		Username:        request.Username,
		Password:        request.Password,
		Age:             request.Age,
		ProfileImageUrl: request.ProfileImageUrl,
	}

	if db := r.postgres.DB.Create(&user).Scan(&userResponse); db.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(db.Error, &pgErr) {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "uni_users_email" {
				return nil, fmt.Errorf("duplicate email")
			}
		}

		return nil, db.Error
	}

	return &userResponse, nil
}

func (r *UserRepo) FindUserById(id int) (*core.User, error) {
	var userResponse core.User
	if db := r.postgres.DB.First(&userResponse, id); db.Error != nil {
		return nil, db.Error
	}

	return &userResponse, nil
}

func (r *UserRepo) FindUserByEmail(email string) (*core.User, error) {
	var userResponse core.User
	db := r.postgres.DB.Where("email = ?", email).First(&userResponse)
	if db.Error != nil {
		return nil, db.Error
	}

	return &userResponse, nil
}

func (r *UserRepo) UpdateUser(user *core.UserUpdateRequest, id int) (*core.UserResponse, error) {
	var userResponse core.UserResponse
	db := r.postgres.DB.Model(&core.User{}).Where("id = ?", id).Updates(&user).First(&userResponse)
	if db.Error != nil {
		return nil, db.Error
	}

	return &userResponse, nil
}

func (r *UserRepo) DeleteUser(id int) error {
	if db := r.postgres.DB.Delete(&core.User{}, id); db.Error != nil {
		return db.Error
	}

	return nil
}
