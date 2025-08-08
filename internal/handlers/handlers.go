package handlers

import (
	"go.uber.org/zap"
	"partiuFit/internal/store"
)

type Handlers struct {
	WorkoutHandlers *WorkoutsHandlers
	UserHandlers    *UserHandlers
	TokensHandlers  *TokensHandlers
	Logger          *zap.SugaredLogger
}

func NewHandlers(store *store.Store, logger *zap.SugaredLogger) *Handlers {
	return &Handlers{
		WorkoutHandlers: NewWorkoutsHandlers(store, logger),
		UserHandlers:    NewUserHandlers(store, logger),
		TokensHandlers:  NewTokensHandlers(store, logger),
		Logger:          logger,
	}
}
