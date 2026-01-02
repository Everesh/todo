package main

import (
	"fmt"
	"os"

	"github.com/nirabyte/todo/internal/app"
)

func main() {
	app := app.New()
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

