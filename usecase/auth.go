package usecase

import "golang.org/x/crypto/bcrypt"

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func verifyPassword(userPassword string, providedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
}
