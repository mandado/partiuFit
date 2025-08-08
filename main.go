package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"net/http"
	"os"
	"partiuFit/internal/app"
	"partiuFit/internal/routes"
	"partiuFit/internal/utils"
	"time"
)

func main() {
	utils.MustIfError(godotenv.Load())
	port := utils.StringToInt(os.Getenv("PORT"))
	logger := zap.Must(zap.NewProduction()).Sugar()

	application, err := app.NewApplication(logger)

	if err != nil {
		panic(err)
	}

	defer application.DB.Close()

	application.Logger.Info("Application started")

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
		Handler:      routes.RegisterRoutes(application),
	}

	application.Logger.Info(fmt.Sprintf("Server started at port %d", port))
	if err := server.ListenAndServe(); err != nil {
		application.Logger.Fatal("Server failed to start", zap.Error(err))
	}
}
