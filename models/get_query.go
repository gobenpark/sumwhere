package models

type GetQuery struct {
	SortBy  []string `json:"sortBy" valid:"optional" example:"email"`
	OrderBy []string `json:"orderBy" valid:"optional" example:"desc"`
	Offset  int      `json:"offset" valid:"numeric,required" example:"1"`
	Limit   int      `json:"limit" valid:"numeric,required" example:"2"`
}
