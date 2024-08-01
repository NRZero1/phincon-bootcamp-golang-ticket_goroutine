package impldb

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"ticket_goroutine/internal/domain"
	"ticket_goroutine/internal/repository"

	"github.com/rs/zerolog/log"
)

type UserRepository struct {
	mtx sync.Mutex
	db *sql.DB
}

func NewUserRepository(database *sql.DB) repository.UserRepositoryInterface {
	return &UserRepository{
		db: database,
	}
}

func (repo *UserRepository) Save(ctx context.Context, user *domain.User) (error) {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	log.Trace().Msg("Inside user repository save")

	select {
	case <- ctx.Done():
		log.Error().Msg(fmt.Sprintf("Error when trying to save user because of timeout with message: %s", ctx.Err()))
		return ctx.Err()
	default:
		log.Trace().Msg("Attempting to save new user")
		
		trx, err := repo.db.BeginTx(ctx, nil)
		log.Trace().Msg("Begin transaction")

		if err != nil {
			return err
		}

		query := `INSERT INTO users (email, name, phone_number, balance) VALUES ($1, $2, $3, $4) RETURNING user_id`
		log.Trace().Msg("Query is set")

		stmt, err := trx.PrepareContext(ctx, query)
		log.Trace().Msg("Prepared statement created with context")

		if err != nil {
			return err
		}

		defer stmt.Close()

		errScan := stmt.QueryRowContext(ctx, user.Email, user.Name, user.PhoneNumber, user.Balance).Scan(&user.UserID)
		log.Trace().Msg("Query ran")

		if errScan != nil {
			trx.Rollback()
			return errScan
		}

		if err = trx.Commit(); err != nil {
			return err
		}

		log.Info().Msg("New user saved")
		return nil
	}
}

func (repo *UserRepository) FindByID(ctx context.Context, id int) (domain.User, error) {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	log.Trace().Msg("Inside user repository find by id")
	log.Debug().Msg(fmt.Sprintf("User repo find by id received id with value %d", id))
	
	select {
	case <- ctx.Done():
		log.Error().Msg(fmt.Sprintf("Error when trying to fetch user because of timeout with message: %s", ctx.Err()))
		return domain.User{}, ctx.Err()
	default:
		log.Trace().Msg("Attempting to fetch user")
		trx, err := repo.db.BeginTx(ctx, nil)
		log.Trace().Msg("Begin transaction")

		if err != nil {
			return domain.User{}, err
		}

		query := "SELECT * FROM users WHERE user_id=$1"
		log.Trace().Msg("Query is set")

		stmt, err := trx.PrepareContext(ctx, query)
		log.Trace().Msg("Prepared statement created with context")

		if err != nil {
			return domain.User{}, err
		}

		defer stmt.Close()

		var user domain.User
		errScan := stmt.QueryRowContext(ctx, id).Scan(
			&user.UserID,
			&user.Email,
			&user.Name,
			&user.PhoneNumber,
			&user.Balance,
		)

		if errScan != nil {
			if errScan == sql.ErrNoRows {
				log.Error().Msg(fmt.Sprintf("User with ID %d not found", id))
				return domain.User{}, fmt.Errorf("no user found with ID %d", id)
			}
			return domain.User{}, errScan
		}

		if err = trx.Commit(); err != nil {
			return domain.User{}, err
		}

		log.Info().Msg("User repo find by id completed")
		return user, nil
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
		trx, err := repo.db.BeginTx(ctx, nil)
		log.Trace().Msg("Begin Transaction")

		if err != nil {
			return []domain.User{}, err
		}

		query := "SELECT * FROM users"
		log.Trace().Msg("Query is set")

		stmt, err := trx.PrepareContext(ctx, query)
		log.Trace().Msg("Prepared statement created with context")

		if err != nil {
			return []domain.User{}, err
		}

		defer stmt.Close()

		res, err := stmt.QueryContext(ctx)
		log.Trace().Msg("Query ran")

		if err != nil {
			return []domain.User{}, err
		}

		defer res.Close()

		var listOfUsers []domain.User

		for res.Next() {
			var user domain.User
			res.Scan(&user.UserID, &user.Email, &user.Name, &user.PhoneNumber, &user.Balance)

			listOfUsers = append(listOfUsers, user)
		}

		log.Info().Msg("Fetching completed")
		return listOfUsers, nil
	}
}

func (repo *UserRepository) ReduceBalance(ctx context.Context, id int, amount float64) (domain.User, error) {
	defer repo.mtx.Unlock()

	log.Trace().Msg("Inside user repository reduce balance")
	log.Debug().Msg(fmt.Sprintf("User repo reduce balance received id with value %d and amount %f", id, amount))

	select {
	case <- ctx.Done():
		log.Error().Msg(fmt.Sprintf("Error when trying to reduce user balance because of timeout with message: %s", ctx.Err()))
		return domain.User{}, ctx.Err()
	default:
		log.Trace().Msg("Attempting to reduce user balance")
		foundUser, err := repo.FindByID(ctx, id)
		log.Trace().Msg("Begin transaction")

		if err != nil {
			return domain.User{}, err
		}

		repo.mtx.Lock()
		log.Debug().Msg(fmt.Sprintf("Balance before reduced: %f", foundUser.Balance))
		foundUser.Balance = foundUser.Balance - amount
		
		trx, err := repo.db.BeginTx(ctx, nil)

		if err != nil {
			return domain.User{}, err
		}

		query := `UPDATE users SET balance=$1 WHERE user_id=$2 RETURNING *`
		log.Trace().Msg("Query is set")

		stmt, err := trx.PrepareContext(ctx, query)
		log.Trace().Msg("Prepared statement created with context")

		if err != nil {
			return domain.User{}, err
		}

		defer stmt.Close()

		var user domain.User

		errScan := stmt.QueryRowContext(ctx, foundUser.Balance, id).Scan(&user.UserID, &user.Email, &user.Name, &user.PhoneNumber, &user.Balance)
		log.Trace().Msg("Query ran")

		if errScan != nil {
			return domain.User{}, errScan
		}

		if err = trx.Commit(); err != nil {
			return domain.User{}, err
		}

		log.Debug().Msg(fmt.Sprintf("Balance after reduced: %f", user.Balance))
		log.Info().Msg("Balance reduced successfully")
		return user, nil
	}
}