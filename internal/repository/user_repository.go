package repository

import (
	"context"
	"ticket_goroutine/internal/domain"
)

type UserRepositoryInterface interface {
	UserSave
	UserFindById
	UserGetAll
	UserBalanceReduce
	UserFindByEmail
}

type UserSave interface {
	Save(context context.Context, user *domain.User) (error)
}

type UserFindById interface {
	FindByID(context context.Context, id int) (domain.User, error)
}

type UserGetAll interface {
	GetAll(context context.Context) ([]domain.User, error)
}

type UserBalanceReduce interface {
	ReduceBalance(context context.Context, id int, amount float64) (domain.User, error)
}

type UserFindByEmail interface {
	FindByEmail(context context.Context, email string) (domain.User, error)
}