package repository

import (
	"context"
	"github.com/liang21/webook/internal/domain"
	"github.com/liang21/webook/internal/repository/dao"
)

type UserRepository struct {
	dao *dao.UserDao
}

func NewUserRepository(dao *dao.UserDao) *UserRepository {
	return &UserRepository{dao: dao}
}

func (u *UserRepository) Create(ctx context.Context, user domain.User) error {
	return u.dao.Insert(ctx, dao.User{
		Email:    user.Email,
		Password: user.Password,
	})
}

func (u *UserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := u.dao.GetByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return u.toDomain(user), nil
}

func (u *UserRepository) toDomain(user dao.User) domain.User {
	return domain.User{
		Id:       user.Id,
		Email:    user.Email,
		Password: user.Password,
	}
}
