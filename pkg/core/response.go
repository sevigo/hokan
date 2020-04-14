package core

type ErrorResp struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
}

type LinksResp struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

type MetaDataResp struct {
	TotalItems int `json:"total_items"`
}

type FilesListResp struct {
	Files []*File      `json:"files"`
	Links []LinksResp  `json:"links"`
	Meta  MetaDataResp `json:"meta"`
}

type DirectoriesListResp struct {
	Directories []*Directory `json:"directories"`
	Links       []LinksResp  `json:"links"`
	Meta        MetaDataResp `json:"meta"`
}

type APIListResp struct {
	Links []LinksResp `json:"links"`
}
