package handlers

import (
	"fmt"
	"net/http"
)

func (h *Handlers) HeathCheck(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "OK")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
