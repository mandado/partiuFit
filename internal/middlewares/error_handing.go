package middlewares

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"net/http"
	internalErrors "partiuFit/internal/errors"
	"partiuFit/internal/utils"
	"runtime/debug"
)

type ErrorHandlerMiddleware struct {
	Logger *zap.SugaredLogger
}

func NewErrorHandlerMiddleware(logger *zap.SugaredLogger) *ErrorHandlerMiddleware {
	return &ErrorHandlerMiddleware{
		Logger: logger,
	}
}

func (em *ErrorHandlerMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if rvr == http.ErrAbortHandler {
					// we don't recover http.ErrAbortHandler so the response
					// to the client is aborted, this should not be logged
					panic(rvr)
				}

				err := rvr.(error)
				validationErrors := &validator.ValidationErrors{}

				logEntry := middleware.GetLogEntry(r)
				if logEntry != nil {
					logEntry.Panic(rvr, debug.Stack())
				} else {
					middleware.PrintPrettyStack(rvr)
				}

				if errors.Is(err, internalErrors.NoRows) {
					em.Logger.Error(err)
					utils.MustWriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "not found"})
					return
				}

				if errors.Is(err, internalErrors.InvalidCredentials) {
					em.Logger.Error(err)
					utils.MustWriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid credentials1"})
					return
				}

				if errors.Is(err, internalErrors.Forbidden) {
					em.Logger.Error(err)
					utils.MustWriteJSON(w, http.StatusForbidden, utils.Envelope{"error": err.Error()})
					return
				}

				if errors.Is(err, internalErrors.InvalidIDParam) {
					em.Logger.Error(err)
					utils.MustWriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
					return
				}

				if errors.Is(err, internalErrors.InvalidIDType) {
					em.Logger.Error(err)
					utils.MustWriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
					return
				}

				if errors.As(err, &validationErrors) {
					validationMap := make(map[string][]string)

					for _, validationError := range *validationErrors {
						fieldName := validationError.Field()
						validationMap[fieldName] = append(validationMap[fieldName], validationError.Translate(utils.Trans))
					}

					em.Logger.Error(err)
					utils.MustWriteJSON(w, http.StatusBadRequest, utils.Envelope{"errors": validationMap})
					return
				}

				if r.Header.Get("Connection") != "Upgrade" {
					w.WriteHeader(http.StatusInternalServerError)
					utils.MustWriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}
