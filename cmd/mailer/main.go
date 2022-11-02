package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
	"os"
	cfg "smtp-client/cmd/mailer/config"
	"smtp-client/internal/api/rest"
	"smtp-client/internal/mailer"
	"smtp-client/internal/mailer/util/html"
	"smtp-client/internal/mailer/util/memq"
	"smtp-client/internal/mailer/util/memsched"
	"smtp-client/internal/mailer/util/plugins/viewstat"
	"smtp-client/internal/mailer/util/plugins/viewstat/redistat"
	"smtp-client/internal/mailer/util/postgrepo"
	"smtp-client/internal/mailer/util/smtp"
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

	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.ViewStat.Redis.Addr,
		Username: config.ViewStat.Redis.Username,
		Password: config.ViewStat.Redis.Password,
		DB:       config.ViewStat.Redis.DB,
	})

	db := postgrepo.New(postgres, html.NewEngine())

	app := mailer.New(
		memsched.New[string](),
		memq.New[string](2048),
		smtp.Delivery{},
		db,
	)

	redisVisitor := redistat.New(
		config.ViewStat.URLs,
		config.ViewStat.Redis.Channel,
		redisClient)
	vs := viewstat.New(redisVisitor)

	app.UsePlugin(vs)

	go redisVisitor.Run(context.TODO())
	go vs.Run(context.TODO())
	go app.Run(context.TODO())

	addr := fmt.Sprintf("%s:%d",
		config.API.HTTP.Host,
		config.API.HTTP.Port)

	rest.New(app).Run(addr)
}
