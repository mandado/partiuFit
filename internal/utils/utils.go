package utils

import (
	"encoding/json"
	"net/http"
	internalErrors "partiuFit/internal/errors"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Envelope map[string]interface{}

func WriteJSON(w http.ResponseWriter, status int, data Envelope) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		return err
	}
	return nil
}

func MustWriteJSON(w http.ResponseWriter, status int, data Envelope) {
	MustIfError(WriteJSON(w, status, data))
}

func MustReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) {
	maxBytes := 1024 * 1024 // 1 MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		panic(err)
	}
}

func ReadIDParam(r *http.Request) (int, error) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		return 0, internalErrors.InvalidIDParam
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return 0, internalErrors.InvalidIDType
	}

	return int(id), nil
}

func Must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}

func MustIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func ValueToPointer[T any](value T) *T {
	return &value
}

func StringToInt(value interface{}) int {
	intValue := Must(strconv.Atoi(value.(string)))
	return intValue
}
