package service

import (
	"context"
	"github.com/liang21/webook/internal/domain"
	"github.com/liang21/webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Login(ctx context.Context, email, password string) (user domain.User, err error)
	SignUp(ctx context.Context, user domain.User) error
	Edit(ctx context.Context, user domain.User) error
	Profile(ctx context.Context, id int64) (domain.User, error)
}

type userService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (u *userService) SignUp(ctx context.Context, user domain.User) error {
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(password)
	return u.repo.Create(ctx, user)
}

func (u *userService) Login(ctx context.Context, email, password string) (user domain.User, err error) {
	user, err = u.repo.GetByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u *userService) Edit(ctx context.Context, user domain.User) error {
	return u.repo.Update(ctx, user)
}

func (u *userService) Profile(ctx context.Context, id int64) (domain.User, error) {
	user, err := u.repo.GetById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}
