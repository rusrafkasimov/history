package main

import (
	"context"
	"flag"
	"fmt"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"time"
)

// @title Swagger Account History Service
// @version 0.1
// @description Account History Microservice (Golang)

// @contact.name Ruslan Kasimov

// @host 127.0.0.1:8091
// @BasePath /

// @securityDefinitions.apikey TokenJWT
// @in header
// @name Authorization

const (
	Name           = "History"
	contextKeyName = "Name"
	serverTimeout  = 10 * time.Second
)

func main() {
	ctx := context.Background()

	id, err := gonanoid.New()
	if err != nil {
		fmt.Println("Can't generate new node ID")
		return
	}

	ctx = context.WithValue(
		ctx,
		contextKeyName,
		Name+"_"+id,
	)

	var env string

	flag.StringVar(&env, "env", ".env.local", "Environment Variables filename")
	flag.Parse()
}
