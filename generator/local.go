package main

import (
	"context"
	"fmt"
)

type localContextKey string

var k localContextKey = localContextKey("localKey")
var ctx context.Context = context.WithValue(context.Background(), k, "LocalValue")

func local() {
	resp, err := Handler(ctx)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp)
}
