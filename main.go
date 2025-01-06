package main

import (
	"fmt"
	"github.com/levensspel/go-gin-template/di"
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/levensspel/go-gin-template/server"
)

func main() {
	health := di.Injector.HealthCheck()
	fmt.Println("DI HealthCheck: %v\n", health)

	err := server.Start()
	if err != nil {
		log.Fatalln(err)
	}

}
