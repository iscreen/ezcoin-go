package response

type PageMeta struct {
	Total     int64 `json:"totalCount"`
	Page      int   `json:"page"`
	PageSize  int   `json:"pageSize"`
	TotalPage int   `json:"totalPage"`
}

type PageResult struct {
	List interface{} `json:"list"`
	Meta PageMeta    `json:"meta"`
}
