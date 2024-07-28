package usecase

import (
	"context"
	"ticket_goroutine/internal/domain"
)

type UserUseCaseInterface interface {
	UserSave
	UserFindById
	UserGetAll
	UserBalanceReduce
}

type UserSave interface {
	Save(context context.Context, user domain.User) (domain.User, error)
}

type UserFindById interface {
	FindById(context context.Context, id int) (domain.User, error)
}

type UserGetAll interface {
	GetAll(context context.Context) ([]domain.User, error)
}

type UserBalanceReduce interface {
	ReduceBalance(context context.Context, id int, amount float64) (domain.User, error)
}