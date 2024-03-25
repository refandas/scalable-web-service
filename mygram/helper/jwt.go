package helper

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"strings"
)

const SECRETKEY = "MawarITUbiRu,vIoletItumErah"

func GenerateToken(id int64, email string) (string, error) {
	claims := jwt.MapClaims{
		"id":    id,
		"email": email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	res, err := token.SignedString([]byte(SECRETKEY))

	return res, err
}

func VerifyToken(ctx *gin.Context) (interface{}, error) {
	err := errors.New("please login to get the token")
	auth := ctx.Request.Header.Get("Authorization")
	bearer := strings.HasPrefix(auth, "Bearer")

	if !bearer {
		return nil, err
	}

	tokenStr := strings.Split(auth, "Bearer ")[1]

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRETKEY), nil
	})

	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(jwt.MapClaims); !ok {
		return nil, err
	}

	return token.Claims.(jwt.MapClaims), nil
}
