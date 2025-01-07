package main

import (
	"fmt"
	"github.com/levensspel/go-gin-template/di"
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/levensspel/go-gin-template/server"
)

func main() {
	healthCheckDI()

	err := server.Start()
	if err != nil {
		log.Fatalln(err)
	}

}

func healthCheckDI() {
	health := di.Injector.HealthCheck()
	fmt.Printf("DI HealthCheck: %v\n", health)
	isHealthy := true
	for service, err := range health {
		if err != nil {
			fmt.Printf("Service %s is unhealthy: %v\n", service, err)
			isHealthy = false
		}
	}
	if !isHealthy {
		panic("DI is not healthy")
	}
}
