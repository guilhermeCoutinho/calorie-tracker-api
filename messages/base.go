package messages

type BaseResponse struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}
