package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var encryptionKey = []byte("659B832C2A8581A8F3429C931A00208E")

const (
	userIDClaim      = "userId"
	accessLevelClaim = "accessLevel"
	expirationClaim  = "exp"
)

func (u *Usecase) Login(ctx context.Context, userName, password string) (string, error) {
	user, err := u.dal.User.GetUser(ctx, userName)
	if err != nil {
		return "", err
	}

	err = verifyPassword(user.Password, password)
	if err != nil {
		return "", err
	}

	token, err := getToken(user.ID, time.Hour*1)
	return token, err
}

func getToken(userID uuid.UUID, expiration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims[accessLevelClaim] = true
	claims[userIDClaim] = userID
	claims[expirationClaim] = time.Now().Add(expiration).Unix()

	token.Claims = claims

	tokenString, err := token.SignedString(encryptionKey)
	if err != nil {
		return "", errors.New("failed to sign token " + err.Error())
	}

	return tokenString, nil
}

func (u *Usecase) UUIDFromToken(tokenStr string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid sign method")
		}
		return encryptionKey, nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	userIDStr := claims[userIDClaim].(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, err
	}
	return userID, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func verifyPassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
