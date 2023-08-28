package models

type User struct {
	Username   string `gorm:"index;type:string;size:255;comment:名字"`
	Password   string `gorm:"type:string;size:255;comment:密码"`
	Lang       string `gorm:"type:string;size:50;comment:语言"`
	GatewayKey string `gorm:"type:string;size:255;comment:密钥"`
	OpenAPIKey string `gorm:"type:string;size:255;comment:密钥"`
}

func (table *User) Get() interface{} {
	return *table
}

func (table *User) TableName() string {
	return "user"
}
