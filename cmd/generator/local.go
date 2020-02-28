package main

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"
)

type localContextKey string

var k localContextKey = localContextKey("localKey")
var ctx context.Context = context.WithValue(context.Background(), k, "LocalValue")

func local() {
	err := godotenv.Load()

	resp, err := Handler(ctx)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp)
}
