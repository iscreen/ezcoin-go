package request

// PageInfo Paging common input parameter structure
type PageInfo struct {
	Page     int `json:"page" form:"page"`         // 頁碼
	PageSize int `json:"pageSize" form:"pageSize"` // 每頁大小
}

// GetById Find by id structure
type GetById struct {
	ID float64 `json:"id" form:"id"` // Primary Key
}

func (r *GetById) Uint() uint {
	return uint(r.ID)
}

type IdsReq struct {
	Ids []int `json:"ids" form:"ids"`
}

// GetAuthorityId Get role by id structure
type GetAuthorityId struct {
	AuthorityId string `json:"authorityId" form:"authorityId"` // 角色 ID
}

type TableQuery struct {
	Page     int               `json:"page" form:"page"`         // 頁碼
	PageSize int               `json:"pageSize" form:"pageSize"` // 每頁大小
	Sort     string            `json:"sort" form:"sort"`
	Order    string            `json:"order" form:"order"`
	Filter   map[string]string `json:"filter" form:"filter"`
}

type Empty struct{}
