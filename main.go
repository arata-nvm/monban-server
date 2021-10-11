package main

import (
	"log"

	"github.com/arata-nvm/monban/database"
	"github.com/arata-nvm/monban/env"
	"github.com/arata-nvm/monban/web"
)

func main() {
	err := database.Initialize()
	if err != nil {
		log.Fatal(err)
	}

	e := web.NewRouter()
	e.Logger.Fatal(e.Start(":" + env.Port()))
}
