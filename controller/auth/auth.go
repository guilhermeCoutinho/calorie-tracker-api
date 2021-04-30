package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/guilhermeCoutinho/api-studies/dal"
	"github.com/guilhermeCoutinho/api-studies/messages"
	"github.com/guilhermeCoutinho/api-studies/models"
	"github.com/guilhermeCoutinho/api-studies/server/http/wrapper"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

var encryptionKey = []byte("659B832C2A8581A8F3429C931A00208E")

const (
	userIDClaim      = "userId"
	accessLevelClaim = "accessLevel"
	expirationClaim  = "exp"
)

type Auth struct {
	dal    *dal.DAL
	config *viper.Viper
}

func NewAuth(
	dal *dal.DAL,
	config *viper.Viper,
) *Auth {
	return &Auth{
		dal:    dal,
		config: config,
	}
}

func (a *Auth) Post(ctx context.Context, args *messages.LoginRequest, vars *struct{}) (*messages.LoginResponse, *wrapper.HandlerError) {
	user, err := a.dal.User.GetUser(ctx, args.Username, nil)
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusNotFound}
	}

	err = verifyPassword(user.Password, args.Password)
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusUnauthorized}
	}

	token, err := getToken(user, time.Hour*1)
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusUnauthorized}
	}

	response := &messages.LoginResponse{
		AccessToken: token,
	}
	return response, nil
}

func getToken(user *models.User, expiration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims[accessLevelClaim] = user.AccessLevel
	claims[userIDClaim] = user.ID
	claims[expirationClaim] = time.Now().UTC().Add(expiration).Unix()

	token.Claims = claims

	tokenString, err := token.SignedString(encryptionKey)
	if err != nil {
		return "", errors.New("failed to sign token " + err.Error())
	}

	return tokenString, nil
}

func (a *Auth) ClaimsFromToken(tokenStr string) (*models.Claims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid sign method")
		}
		return encryptionKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if err = claims.Valid(); err != nil {
		return nil, err
	}

	userIDStr := claims[userIDClaim].(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	accessLevel := int(claims[accessLevelClaim].(float64))

	return &models.Claims{
		UserID:      userID,
		AccessLevel: models.AccessLevel(accessLevel),
	}, nil
}

func verifyPassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
