package handlers

import (
	"net/http"
	"partiuFit/internal/requests"
	"partiuFit/internal/store"
	"partiuFit/internal/tokens"
	"partiuFit/internal/utils"
	"time"

	"go.uber.org/zap"
)

type TokensHandlers struct {
	Store  *store.Store
	Logger *zap.SugaredLogger
}

func NewTokensHandlers(store *store.Store, logger *zap.SugaredLogger) *TokensHandlers {
	return &TokensHandlers{
		Store:  store,
		Logger: logger,
	}
}

func (th *TokensHandlers) CreateToken(w http.ResponseWriter, r *http.Request) {
	var request requests.CreateTokenRequest
	utils.MustReadJSON(w, r, &request)
	utils.MustValidateStruct(request)

	user, err := th.Store.UserStore.GetUserByUsername(request.Username)

	if err != nil || user == nil {
		th.Logger.Errorf("GetUserByUsername: failed to get user: %v", err)
		utils.MustWriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid credentials"})
		return
	}

	passwordMatch, err := user.Password.VerifyPassword(request.Password)

	if err != nil {
		th.Logger.Errorf("CheckPassword: failed to check password: %v", err)
		utils.MustWriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if !passwordMatch {
		th.Logger.Error("invalid credentials")
		utils.MustWriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid credentials"})
		return
	}

	expiresAt := time.Hour * 24
	token, err := th.Store.TokensStore.CreateToken(user.ID, expiresAt, tokens.ScopeAuthentication)

	if err != nil {
		th.Logger.Errorf("CreateToken: failed to create token: %v", err)
		utils.MustWriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.MustWriteJSON(w, http.StatusOK, utils.Envelope{"token": token.Plaintext})
}
