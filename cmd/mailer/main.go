package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
	"gopkg.in/yaml.v2"
	"os"
	"smtp-client/internal/api/rest"
	"smtp-client/internal/mailer"
	"smtp-client/internal/mailer/util/html"
	"smtp-client/internal/mailer/util/memq"
	"smtp-client/internal/mailer/util/plugins/viewstat"
	"smtp-client/internal/mailer/util/plugins/viewstat/redistat"
	"smtp-client/internal/mailer/util/postgrepo"
	"smtp-client/internal/mailer/util/smtp"
	"strings"
)

type Config struct {
	DB struct {
		Postgres struct {
			Host     string `json:"host,omitempty" yaml:"host"`
			Port     int    `json:"port,omitempty" yaml:"port"`
			DB       string `json:"db,omitempty" yaml:"db"`
			User     string `json:"user,omitempty" yaml:"user"`
			Password string `json:"password,omitempty" yaml:"password"`
		} `json:"postgres" yaml:"postgres"`
	} `json:"db" yaml:"db"`

	ViewStat struct {
		Redis struct {
			Channel  string `json:"channel,omitempty" yaml:"channel"`
			Addr     string `json:"addr,omitempty" yaml:"addr"`
			Username string `json:"username,omitempty" yaml:"username"`
			Password string `json:"password,omitempty" yaml:"password"`
			DB       int    `json:"db,omitempty" yaml:"db"`
		} `json:"redis" yaml:"redis"`

		URLs []string `json:"urls,omitempty" yaml:"urls"`
	} `json:"view_stat" yaml:"view_stat"`

	API struct {
		HTTP struct {
			Host string `json:"host,omitempty" yaml:"host"`
			Port int    `json:"port,omitempty" yaml:"port"`
		} `json:"http" yaml:"http"`
	} `json:"api" yaml:"api"`
}

func getConfig(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var config Config

	switch {
	case strings.HasSuffix(path, ".yml"), strings.HasSuffix(path, ".yaml"):
		if err := yaml.NewDecoder(file).Decode(&config); err != nil {
			return Config{}, err
		}
	case strings.HasSuffix(path, ".json"):
		if err := json.NewDecoder(file).Decode(&config); err != nil {
			return Config{}, err
		}
	default:
		return Config{}, errors.New("unsupported format")
	}
	return config, nil
}

func main() {
	configPath := "config.json"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	config, err := getConfig(configPath)
	if err != nil {
		panic(fmt.Sprintf("failed to get the config: %s", err))
	}

	// Postgres
	//postgresLink := "postgresql://localhost:5432/postgres?user=postgres&password=postgres"
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

	app := mailer.New(nil,
		memq.New[string](2048),
		smtp.Delivery{},
		db,
	)

	// 		"",

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
