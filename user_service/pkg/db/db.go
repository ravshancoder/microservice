package db

import (
	"fmt"

	"najottalim/6_part_microservice/service/user_service/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //postgres drivers
)

func ConnectToDB(cfg config.Config) (*sqlx.DB, error) {
	psqlString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDatabase,
	)

	connDb, err := sqlx.Open("postgres", psqlString)
	if err != nil {
		return nil, err
	}

	return connDb, nil
}
