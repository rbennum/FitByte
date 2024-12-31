package server

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/levensspel/go-gin-template/config"
	"github.com/levensspel/go-gin-template/helper"
	"github.com/levensspel/go-gin-template/middleware"
)

func Start() error {
	db, err := config.NewDbInit()

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	helper.WORK_DIR = wd

	r := gin.Default()
	r.Use(middleware.EnableCORS)

	NewRouter(r, db)

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

		err := r.RunTLS(
			fmt.Sprintf("0.0.0.0:%s", port),
			sslCert,
			sslKey,
		)
		if err != nil {
			log.Fatalf("Failed to start HTTPS server: %v", err)
		}
	default:
		gin.SetMode(gin.DebugMode)
		r.Run(fmt.Sprintf("0.0.0.0:%s", port))
	}

	return nil
}
