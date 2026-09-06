package server

import (
	"errors"
	"net/http"

	"github.com/RA9/gamifydev/platform/internal/store"
)

func writeWorkAccessError(w http.ResponseWriter, err error) bool {
	switch {
	case errors.Is(err, store.ErrWorkLocked):
		http.Error(w, "This practical work has not opened yet. Follow your cohort schedule and complete the earlier required work first.", http.StatusForbidden)
		return true
	case errors.Is(err, store.ErrAccountRestricted):
		http.Error(w, "This account cannot participate right now.", http.StatusForbidden)
		return true
	case err != nil:
		http.Error(w, "could not verify access to this work", http.StatusInternalServerError)
		return true
	default:
		return false
	}
}
