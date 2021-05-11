package core

type Resp struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
	Msg    string `json:"message"`
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

type TargetConfigsListResp struct {
	Links []LinksResp  `json:"links"`
	Meta  MetaDataResp `json:"meta"`
}

type TargetListResp struct {
	Links []LinksResp  `json:"links"`
	Meta  MetaDataResp `json:"meta"`
}

type DirectoriesListResp struct {
	Directories []*Directory `json:"directories"`
	Links       []LinksResp  `json:"links"`
	Meta        MetaDataResp `json:"meta"`
}

type DirectoryDetails struct {
	Directory Directory      `json:"directory"`
	Stats     DirectoryStats `json:"stats"`
	Links     []LinksResp    `json:"links"`
}

type APIListResp struct {
	Links []LinksResp `json:"links"`
}
