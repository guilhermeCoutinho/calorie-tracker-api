package messages

type LoginRequest struct {
	Username string `json:"user_name"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"token"`
}
