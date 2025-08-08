package app

import (
	"database/sql"
	"os"
	"partiuFit/internal/database"
	"partiuFit/internal/handlers"
	"partiuFit/internal/middlewares"
	"partiuFit/internal/store"
	"partiuFit/migrations"

	"go.uber.org/zap"
)

type Middlewares struct {
	UserMiddleware         *middlewares.UserMiddleware
	ErrorHandlerMiddleware *middlewares.ErrorHandlerMiddleware
	SecurityMiddleware     *middlewares.SecurityMiddleware
}

type Application struct {
	Logger      *zap.SugaredLogger
	Handlers    *handlers.Handlers
	Middlewares Middlewares
	DB          *sql.DB
}

func NewApplication(logger *zap.SugaredLogger) (*Application, error) {
	logger.Info("Initializing application")

	logger.Info("Opening database")
	db, err := database.Open(os.Getenv("DATABASE_URL"))

	if err != nil {
		logger.Fatal("failed to open database", zap.Error(err))
	}

	logger.Info("Migrating database")
	err = database.MigrateFS(db, migrations.FS, migrations.FSPath)

	if err != nil {
		panic(err)
	}

	appStore := store.NewStore(db)
	appHandlers := handlers.NewHandlers(appStore, logger)
	userMiddleware := middlewares.NewUserMiddleware(appStore, logger)
	errorHandlerMiddleware := middlewares.NewErrorHandlerMiddleware(logger)
	securityMiddleware := middlewares.NewSecurityMiddleware(logger)

	app := &Application{
		Logger:   logger,
		Handlers: appHandlers,
		DB:       db,
		Middlewares: Middlewares{
			UserMiddleware:         userMiddleware,
			ErrorHandlerMiddleware: errorHandlerMiddleware,
			SecurityMiddleware:     securityMiddleware,
		},
	}

	return app, nil
}
