package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/levensspel/go-gin-template/cache"
	"github.com/levensspel/go-gin-template/config"
	"github.com/levensspel/go-gin-template/di"
	"github.com/levensspel/go-gin-template/migration"

	_ "github.com/joho/godotenv/autoload"
	"github.com/levensspel/go-gin-template/server"
)

func main() {
	healthCheckDI()

	if config.EnableAutoMigrate() {
		migration.AutoMigrate()
	}

	cache.Initialize()
	// Handle graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	defer cache.Cache.Close()

	go func() {
		<-sig
		cache.Cache.Close()
		os.Exit(0)
	}()

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
