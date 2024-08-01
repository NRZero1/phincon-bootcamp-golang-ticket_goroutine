package db

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
)

type Database struct {
	Db string
}

func NewConnection(database string) Database {
	return Database{
		Db: database,
	}
}

func (database Database) GetConnection(username string, password string, host string, port string, databaseName string) (*sql.DB, error) {
	var driver string
	var connString string

	if database.Db == "mysql" {
		driver = "mysql"

		// format : "username:password@tcp(host:port)/database_name"
		connString = fmt.Sprintf("%s:%s@tcp(%s%s)/%s", username, password, host, port, databaseName)
		log.Debug().Msg(fmt.Sprintf("Conn String is set as: %s", connString))
	} else if database.Db == "postgresql" {
		driver = "pgx"

		// format : "postgres://username:password@localhost:5432/database_name"
		connString = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, databaseName)
		log.Debug().Msg(fmt.Sprintf("Conn String is set as: %s", connString))
	} else {
		return nil, fmt.Errorf("unsupported database detected")
	}

	log.Debug().Str("Driver: ", driver).Msg("Driver value")
	db, err := sql.Open(driver, connString)

	if err != nil {
		log.Error().Msg(fmt.Sprintf("unable to connect to database %s", database.Db))
		return nil, err
	}

	log.Info().Msg(fmt.Sprintf("Running %s on %s on port %s", database.Db, host, port))

	return db, nil
}