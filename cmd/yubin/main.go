package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
	"os"
	cfg "yubin/cmd/yubin/config"
	"yubin/lib/api/rest"
	"yubin/lib/impl/common/codec/jsoncodec"
	redikv "yubin/lib/impl/common/data/kv/redis"
	mongorec "yubin/lib/impl/common/data/record/mongo"
	"yubin/lib/impl/smtp"
	postgres2 "yubin/lib/impl/templates/postgres"
	"yubin/lib/plugins/viewstat"
	"yubin/lib/plugins/viewstat/redistat"
	yubin "yubin/src"
	"yubin/src/user"
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

	redisVisitor := redistat.New(
		config.ViewStat.URLs,
		config.ViewStat.Redis.Channel,
		redisClient)
	vs := viewstat.New(redisVisitor)

	app := yubin.Configure().
		Transport(smtp.Delivery{}).
		Plugins(vs).
		Sources(redikv.New[yubin.NamedSource](redisClient, jsoncodec.New[yubin.NamedSource](), "SOURCES_", "SOURCE_KEYS")).
		Publications(redikv.New[user.User](redisClient, jsoncodec.New[user.User](), "PUBLICATIONS_", "PUBLICATION_KEYS")).
		Users(redikv.New[user.User](redisClient, jsoncodec.New[user.User](), "USERS_", "USER_KEYS")).
		Reports(mongorec.New()).
		Templates(postgres2.Storage{})

	//.New(
	//	memsched.New[string](),
	//	memq.New[string](2048),
	//	smtp.Delivery{},
	//	db,
	//)

	app.UsePlugin(vs)

	go redisVisitor.Run(context.TODO())
	go vs.Run(context.TODO())
	go app.Run(context.TODO())

	addr := fmt.Sprintf("%s:%d",
		config.API.HTTP.Host,
		config.API.HTTP.Port)

	rest.New(app).Run(addr)
}
