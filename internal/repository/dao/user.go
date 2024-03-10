package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var ErrDuplicateEmail = errors.New("邮箱冲突")

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}

func (u *UserDao) Insert(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	user.CreatedAt = now
	user.UpdatedAt = now
	err := u.db.WithContext(ctx).Create(&user).Error
	var me *mysql.MySQLError
	if errors.As(err, &me) {
		if me.Number == 1062 {
			return ErrDuplicateEmail
		}
	}
	return err
}

func (u *UserDao) GetByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := u.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return user, err

}

type User struct {
	Id        int64  `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func (*User) TableName() string {
	return "user"
}
