package main

import (
	"fmt"
	"os"
)

func HandleError(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
