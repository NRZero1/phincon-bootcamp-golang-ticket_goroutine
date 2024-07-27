package repository

import (
	"context"
	"ticket_goroutine/internal/domain"
)

type UserRepositoryInterface interface {
	UserSave
	UserFindById
	UserGetAll
}

type UserSave interface {
	Save(context context.Context, user *domain.User) (domain.User, error)
}

type UserFindById interface {
	FindByID(context context.Context, id int) (domain.User, error)
}

type UserGetAll interface {
	GetAll(context context.Context) ([]domain.User, error)
}