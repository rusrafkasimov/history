package main

import (
	"context"
	"flag"
	"fmt"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"time"
)

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
