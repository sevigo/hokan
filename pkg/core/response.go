package core

type ErrorResp struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
}

type LinksResp struct {
	Href   string `json:"links"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}
