package impl

import (
	"ticket_goroutine/internal/domain"
	"ticket_goroutine/internal/repository"
	"ticket_goroutine/internal/usecase"
)

type EventUseCase struct {
	repo repository.EventRepositoryInterface
}

func NewEventUseCase(repo repository.EventRepositoryInterface) (usecase.EventUseCaseInterface) {
	return EventUseCase{
		repo: repo,
	}
}

func (uc EventUseCase) Save(event domain.Event) (domain.Event, error) {
	return uc.repo.Save(&event)
}

func (uc EventUseCase) FindById(id int) (domain.Event, error) {
	return uc.repo.FindByID(id)
}

func (uc EventUseCase) GetAll() ([]domain.Event) {
	return uc.repo.GetAll()
}