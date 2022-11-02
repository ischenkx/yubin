#!/bin/bash

GOOS=linux go build -o bin/mailer cmd/mailer/main.go
GOOS=linux go build -o bin/setup cmd/mailer/setup/main.go

