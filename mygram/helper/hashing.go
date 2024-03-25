package helper

import "golang.org/x/crypto/bcrypt"

func HashPassword(pass string) string {
	salt := 0
	password := []byte(pass)

	hashed, _ := bcrypt.GenerateFromPassword(password, salt)

	return string(hashed)
}

func ValidateHashPassword(pass, hashedPassword string) bool {
	res := true

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(pass))
	if err != nil {
		res = false
	}

	return res
}
