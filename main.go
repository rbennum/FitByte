package main

import (
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/levensspel/go-gin-template/server"
)

func main() {
	err := server.Start()
	if err != nil {
		log.Fatalln(err)
	}

}
