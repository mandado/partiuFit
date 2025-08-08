package handlers

import (
	"errors"
	"net/http"
	internalErrors "partiuFit/internal/errors"
	"partiuFit/internal/middlewares"
	"partiuFit/internal/requests"
	"partiuFit/internal/store"
	"partiuFit/internal/utils"

	"go.uber.org/zap"
)

type UserHandlers struct {
	Store  *store.Store
	Logger *zap.SugaredLogger
}

func NewUserHandlers(store *store.Store, logger *zap.SugaredLogger) *UserHandlers {
	return &UserHandlers{
		Store:  store,
		Logger: logger,
	}
}

func (uh *UserHandlers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user store.User
	userRequest := &requests.UserRequest{}

	utils.MustReadJSON(w, r, userRequest)
	utils.MustValidateStruct(userRequest)
	utils.MustIfError(user.FromUserRequest(userRequest).HashPassword())

	uh.Logger.Info("creating user", zap.String("name", userRequest.Name))
	err := uh.Store.UserStore.CreateUser(&user)

	if err != nil {
		if errors.Is(err, internalErrors.ErrUserAlreadyExists) {
			uh.Logger.Error("user already exists")
			utils.MustWriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
			return
		}

		uh.Logger.Errorf("failed to create user: %v", err)
		utils.MustWriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create user"})
		return
	}

	utils.MustWriteJSON(w, http.StatusCreated, utils.Envelope{"user": user})
}

func (uh *UserHandlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := middlewares.GetUser(r)
	userRequest := &requests.UserRequest{}

	utils.MustReadJSON(w, r, userRequest)
	user.FromUserRequest(userRequest)

	uh.Logger.Info("updating user", zap.String("name", userRequest.Name))
	err := uh.Store.UserStore.UpdateUser(user.ID, user)

	if err != nil {
		uh.Logger.Errorf("failed to update user: %v", err)
		utils.MustWriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to update user"})
		return
	}

	utils.MustWriteJSON(w, http.StatusOK, utils.Envelope{"user": user})
}
