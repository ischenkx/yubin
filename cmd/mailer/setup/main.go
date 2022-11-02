package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"os"
	cfg "smtp-client/cmd/mailer/config"
	"smtp-client/internal/mailer/util/html"
	"smtp-client/internal/mailer/util/postgrepo"
)

func main() {
	configPath := "config.json"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	config, err := cfg.Read(configPath)
	if err != nil {
		panic(fmt.Sprintf("failed to get the config: %s", err))
	}

	postgresLink := fmt.Sprintf("postgresql://%s:%d/%s?user=%s&password=%s",
		config.DB.Postgres.Host,
		config.DB.Postgres.Port,
		config.DB.Postgres.DB,
		config.DB.Postgres.User,
		config.DB.Postgres.Password,
	)
	postgres, err := pgx.Connect(context.Background(), postgresLink)

	if err != nil {
		panic(fmt.Sprintf("failed to connect to postgres: %s", err))
	}

	db := postgrepo.New(postgres, html.NewEngine())

	db.Init(context.Background())
}
