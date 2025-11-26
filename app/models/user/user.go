package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"column:username"`
	Eamil    string `gorm:"column:mail"`
	Mobile   string `gorm:"column:mobile"`
}

func (u *User) TableName() string {
	return "user_table"
}
