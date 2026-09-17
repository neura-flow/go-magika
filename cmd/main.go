package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	if err := cli(context.Background(), os.Stdout, os.Args[1:]...); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
