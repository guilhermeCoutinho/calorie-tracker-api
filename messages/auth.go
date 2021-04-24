package messages

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	BaseResponse
	AccessToken string `json:"token"`
}
