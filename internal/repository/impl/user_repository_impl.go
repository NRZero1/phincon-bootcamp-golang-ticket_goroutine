package impl

import (
	"context"
	"fmt"
	"sync"
	"ticket_goroutine/internal/domain"
	"ticket_goroutine/internal/repository"

	"github.com/rs/zerolog/log"
)

type UserRepository struct {
	mtx sync.Mutex
	users map[int]domain.User
}

func NewUserRepository() repository.UserRepositoryInterface {
	repo := &UserRepository {
		users: map[int]domain.User{},
	}
	repo.initData()
	return repo
}

func (repo *UserRepository) initData() {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	repo.users[1] = domain.User{
		UserID: 1,
		Email: "ronaldo.pangestu1@gmail.com",
		Name: "Daniel",
		PhoneNumber: "+6282260394746",
		Balance: 1000000,
	}

	repo.users[2] = domain.User{
		UserID: 2,
		Email: "testemail1@gmail.com",
		Name: "test user",
		PhoneNumber: "+1234567890122",
		Balance: 1000000,
	}
}

func (repo *UserRepository) Save(ctx context.Context, user *domain.User) (domain.User, error) {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	log.Trace().Msg("Inside user repository save")

	select {
	case <- ctx.Done():
		log.Error().Msg(fmt.Sprintf("Error when trying to save user because of timeout with message: %s", ctx.Err()))
		return domain.User{}, ctx.Err()
	default:
		log.Trace().Msg("Attempting to save new user")
		if _, exists := repo.users[user.UserID]; exists {
			log.Error().Msg(fmt.Sprintf("User with ID %d already exist", user.UserID))
			return domain.User{}, fmt.Errorf("user with ID %d already exist", user.UserID)
		}

		repo.users[user.UserID] = *user
		log.Info().Msg("New user saved")
		return repo.users[user.UserID], nil
	}
}

func (repo *UserRepository) FindByID(ctx context.Context, id int) (domain.User, error) {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	log.Trace().Msg("Inside user repository find by id")
	
	select {
	case <- ctx.Done():
		log.Error().Msg(fmt.Sprintf("Error when trying to fetch user because of timeout with message: %s", ctx.Err()))
		return domain.User{}, ctx.Err()
	default:
		log.Trace().Msg("Attempting to fetch user")
		if foundUser, exists := repo.users[id]; exists {
			log.Info().Msg("Fetching completed")
			return foundUser, nil
		}
		log.Error().Msg(fmt.Sprintf("User with ID %d not found", id))
		return domain.User{}, fmt.Errorf("no user found with ID %d", id)
	}
}

func (repo *UserRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	log.Trace().Msg("Inside user repository get all")

	select {
	case <- ctx.Done():
		log.Error().Msg(fmt.Sprintf("Error when trying to fetch user because of timeout with message: %s", ctx.Err()))
		return []domain.User{}, ctx.Err()
	default:
		log.Trace().Msg("Attempting to fetch user")
		listOfUsers := make([]domain.User, 0, len(repo.users))

		for _, v := range repo.users {
			listOfUsers = append(listOfUsers, v)
		}

		log.Info().Msg("Fetching completed")
		return listOfUsers, nil
	}
}

func (repo *UserRepository) ReduceBalance(ctx context.Context, id int, amount float64) (domain.User, error) {
	defer repo.mtx.Unlock()

	log.Trace().Msg("Inside user repository reduce balance")

	select {
	case <- ctx.Done():
		log.Error().Msg(fmt.Sprintf("Error when trying to reduce user balance because of timeout with message: %s", ctx.Err()))
		return domain.User{}, ctx.Err()
	default:
		log.Trace().Msg("Attempting to reduce user balance")
		foundUser, err := repo.FindByID(ctx, id)

		if err != nil {
			return domain.User{}, err
		}

		repo.mtx.Lock()
		log.Debug().Msg(fmt.Sprintf("Balance before reduced: %f", foundUser.Balance))
		foundUser.Balance = foundUser.Balance - amount
		repo.users[foundUser.UserID] = foundUser
		log.Debug().Msg(fmt.Sprintf("Balance after reduced: %f", foundUser.Balance))
		log.Info().Msg("Balance reduced successfully")
		return repo.users[foundUser.UserID], nil
	}
}