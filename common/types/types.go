package types

import "time"

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

type GetDataReq struct {
	Tablename string `path:"tablename"`
	Page      int    `form:"page,default=1"`
	Size      int    `form:"size,default=20"`
	Search    string `form:"search,optional"`
	Filter    string `form:"filter,optional"`
	Range     string `form:"range,optional"`
	Sort      string `form:"sort,optional"`
}

type ApplicationTemplate struct {
	Id        int       `yaml:"id" gorm:"primaryKey,autoIncrement"`
	Name      string    `yaml:"name" gorm:"unique"`
	Content   string    `yaml:"content"`
	Variables string    `yaml:"variables"`
	UpdateAt  time.Time `yaml:"update_at"`
}

type User struct {
	Id       int    `gorm:"primaryKey,autoIncrement"`
	Username string `gorm:"size:50;unique"`
	Password string `gorm:"size:100"`
	UpdateAt time.Time
}

type LoginRes struct {
	AccessToken  string `json:"access_token"`
	AccessExpire int64  `json:"access_expire"`
	RefreshAfter int64  `json:"refresh_after"`
}
