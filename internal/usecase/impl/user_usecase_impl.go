package impl

import (
	"context"
	"ticket_goroutine/internal/domain"
	"ticket_goroutine/internal/repository"
	"ticket_goroutine/internal/usecase"

	"github.com/rs/zerolog/log"
)

type UserUseCase struct {
	repo repository.UserRepositoryInterface
}

func NewUserUseCase(repo repository.UserRepositoryInterface) (usecase.UserUseCaseInterface) {
	return UserUseCase{
		repo: repo,
	}
}

func (uc UserUseCase) Save(ctx context.Context, user domain.User) (domain.User, error) {
	log.Trace().Msg("Entering user usecase save")
	return uc.repo.Save(ctx, &user)
}

func (uc UserUseCase) FindById(ctx context.Context, id int) (domain.User, error) {
	log.Trace().Msg("Entering user usecase find by id")
	return uc.repo.FindByID(ctx, id)
}

func (uc UserUseCase) GetAll(ctx context.Context) ([]domain.User, error) {
	log.Trace().Msg("Entering user usecase get all")
	return uc.repo.GetAll(ctx)
}

func (uc UserUseCase) ReduceBalance(ctx context.Context, id int, amount float64) (domain.User, error) {
	log.Trace().Msg("Entering user usecase reduce balance")
	return uc.repo.ReduceBalance(ctx, id, amount)
}