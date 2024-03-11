package repository

import (
	"context"
	"github.com/liang21/webook/internal/domain"
	"github.com/liang21/webook/internal/repository/dao"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	Update(ctx context.Context, user domain.User) error
	GetById(ctx context.Context, id int64) (domain.User, error)
	Create(ctx context.Context, user domain.User) error
}

type CacheUserRepository struct {
	dao *dao.UserDao
}

func NewUserRepository(dao *dao.UserDao) UserRepository {
	return &CacheUserRepository{dao: dao}
}

func (u *CacheUserRepository) Create(ctx context.Context, user domain.User) error {
	return u.dao.Insert(ctx, dao.User{
		Email:    user.Email,
		Password: user.Password,
	})
}

func (u *CacheUserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := u.dao.GetByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return u.toDomain(user), nil
}

func (u *CacheUserRepository) toDomain(user dao.User) domain.User {
	return domain.User{
		Id:       user.Id,
		Email:    user.Email,
		Password: user.Password,
	}
}

func (u *CacheUserRepository) toEntity(user domain.User) dao.User {
	return dao.User{
		Id:       user.Id,
		NikeName: user.NikeName,
		Birthday: user.Birthday.UnixMilli(),
		About:    user.About,
	}
}

func (u *CacheUserRepository) Update(ctx context.Context, user domain.User) error {
	return u.dao.Update(ctx, u.toEntity(user))
}

func (u *CacheUserRepository) GetById(ctx context.Context, id int64) (domain.User, error) {
	user, err := u.dao.GetById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return u.toDomain(user), nil
}
