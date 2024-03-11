package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var ErrDuplicateEmail = errors.New("邮箱冲突")

type UserDao interface {
	Insert(ctx context.Context, user User) error
	GetByEmail(ctx context.Context, email string) (User, error)
	Update(ctx context.Context, user User) error
	GetById(ctx context.Context, id int64) (User, error)
}

type GORMUserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) UserDao {
	return &GORMUserDao{db: db}
}

func (u *GORMUserDao) Insert(ctx context.Context, user User) error {
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

func (u *GORMUserDao) GetByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := u.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return user, err

}

func (u *GORMUserDao) Update(ctx context.Context, user User) error {
	return u.db.WithContext(ctx).Model(&user).Where("id = ?", user.Id).Updates(&user).Error
}

func (u *GORMUserDao) GetById(ctx context.Context, id int64) (User, error) {
	var user User
	err := u.db.WithContext(ctx).First(&user, id).Error
	return user, err
}

type User struct {
	Id        int64  `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	NikeName  string `json:"nike_name"`
	Birthday  int64  `json:"birthday"`
	About     string `json:"about"`
}

func (*User) TableName() string {
	return "user"
}
