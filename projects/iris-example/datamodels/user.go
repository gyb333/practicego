package datamodels

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID             int64     `json:"id" form:"id"`
	Firstname      string    `json:"firstname" form:"firstname"`
	Username       string    `json:"username" form:"username"`
	HashedPassword []byte    `json:"-" form:"-"`
	CreatedAt      time.Time `json:"created_at" form:"created_at"`
}
// IsValid可以做一些非常简单的“低级”数据验证
func (u User) IsValid() bool {
	return u.ID > 0
}
// GeneratePassword将根据我们为我们生成哈希密码
//用户的输入
func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}
// ValidatePassword将检查密码是否匹配
func ValidatePassword(userPassword string, hashed []byte) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(hashed, []byte(userPassword)); err != nil {
		return false, err
	}
	return true, nil
}
