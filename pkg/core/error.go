package core

type ErrorResp struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
}
