package main

import (
	"github.com/arjunmoola/go-arxiv/internal/app"
	"log"
)

func main() {
	app := app.New()

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

