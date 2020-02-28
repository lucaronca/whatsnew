package main

import (
	"fmt"
)

func local() {
	resp, err := Handler(Request{})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp)
}
