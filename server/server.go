package server

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/levensspel/go-gin-template/config"
	dbcontext "github.com/levensspel/go-gin-template/database"
	"github.com/levensspel/go-gin-template/helper"
	"github.com/levensspel/go-gin-template/middleware"
)

func Start() error {
	config := config.LoadConfig()

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	helper.WORK_DIR = wd

	r := gin.Default()
	r.Use(middleware.EnableCORS)

	NewRouter(r, dbcontext.Connect(config.DatabaseURL))

	r.Use(gin.Recovery())

	port := os.Getenv("PORT")

	if len(port) == 0 {
		port = "8080"
	}

	appEnv := os.Getenv("MODE")

	switch appEnv {
	case "PRODUCTION":
		gin.SetMode(gin.ReleaseMode)

		sslCert := os.Getenv("SSL_CERT_PATH")
		sslKey := os.Getenv("SSL_KEY_PATH")

		if sslCert == "" || sslKey == "" {
			log.Fatal("SSL certificates not configured")
		}

		host := os.Getenv("PROD_HOST")
		err := r.RunTLS(
			fmt.Sprintf("%s:%s", host, port),
			sslCert,
			sslKey,
		)
		if err != nil {
			log.Fatalf("Failed to start HTTPS server: %v", err)
		}
	default:
		gin.SetMode(gin.DebugMode)
		host := os.Getenv("DEBUG_HOST")
		r.Run(fmt.Sprintf("%s:%s", host, port))
	}

	return nil
}
