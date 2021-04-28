package test

import (
	"net/http"
	"testing"

	"github.com/guilhermeCoutinho/api-studies/messages"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserSuccess(t *testing.T) {
	drop := GetPG(t)
	drop()

	userName := "userName"
	password := "MyPassword"

	createUserRequest := &messages.CreateUserRequest{
		Username:     userName,
		Password:     password,
		CalorieLimit: 100,
	}

	doRequest(t, http.MethodPost, "/users", nil, createUserRequest, &messages.BaseResponse{})

	loginResponse := &messages.LoginResponse{}
	doRequest(t, http.MethodPost, "/auth", nil, &messages.LoginRequest{
		Username: userName,
		Password: password,
	}, loginResponse)

	assert.NotEmpty(t, loginResponse.AccessToken)

	getUsersResponse := &messages.GetUsersResponse{}
	statusCode := doRequest(t, http.MethodGet, "/users/me", &loginResponse.AccessToken, &struct{}{}, getUsersResponse)

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, 100, getUsersResponse.Users.CalorieLimit)
	assert.Equal(t, userName, getUsersResponse.Users.UserName)
}
