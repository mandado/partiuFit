package middlewares

import (
	"context"
	"net/http"
	"partiuFit/internal/store"
	"partiuFit/internal/tokens"
	"partiuFit/internal/utils"
	"strings"

	"go.uber.org/zap"
)

type contextKey string

type UserMiddleware struct {
	Store  *store.Store
	Logger *zap.SugaredLogger
}

const (
	UserContextKey contextKey = contextKey("user")
)

func NewUserMiddleware(store *store.Store, logger *zap.SugaredLogger) *UserMiddleware {
	return &UserMiddleware{
		Store:  store,
		Logger: logger,
	}
}

func SetUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKey, user)

	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *store.User {
	user, ok := r.Context().Value(UserContextKey).(*store.User)

	if !ok {
		panic("missing user in context")
	}

	return user
}

func (um *UserMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			r = SetUser(r, store.AnonymousUser)

			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.MustWriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid authorization header"})
		}

		token := headerParts[1]
		user, err := um.Store.UserStore.GetUserFromToken(tokens.ScopeAuthentication, token)

		if err != nil {
			utils.MustWriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid token"})
			return
		}

		if user == nil {
			utils.MustWriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "token expired or invalid"})
			return
		}

		r = SetUser(r, user)

		next.ServeHTTP(w, r)
	})
}

func (um *UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)

		if user.IsAnonymous() {
			utils.MustWriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "you must be logged in"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
