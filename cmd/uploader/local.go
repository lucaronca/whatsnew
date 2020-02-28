package main

import (
	"fmt"

	"github.com/joho/godotenv"
)

func local() {
	err := godotenv.Load()

	resp, err := Handler(Request{})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp)
}
