package handler

import (
	"encoding/json"
	"net/http"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode("Hello, World!")
	handleWritingErr(err)
}
