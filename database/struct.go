package database

type Resource struct {
	Code string `json:"code" db:"code,omitupdate"`
	CreateTime string `json:"createTime" db:"createTime,omitcreate,omitupdate"`
	UpdateTime string `json:"updateTime" db:"updateTime,omitcreate,omitupdate"`
	CreatorCode string `json:"creatorCode" db:"creatorCode,omitupdate"`
	CreatorName string `json:"creatorName" db:"creatorName,omitupdate"`
}
