package vo

type PageInfo struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	//Total    int64 `json:"total"`
}

type ListData struct {
	PageInfo PageInfo    `json:"page_info"`
	List     interface{} `json:"list"`
}
