package api

import (
	log "10.254.188.33/matyspi5/pmc/src/logger"
	"encoding/json"
	"fmt"
	"net/http"
)

type apiHandler string

func (h apiHandler) baseHandler(w http.ResponseWriter, r *http.Request) {

	msg := "{\"message\":\"Hello from PMC\"}"
	fmt.Println("Hello from PMC")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(msg)
	if err != nil {
		log.Error("[HANDLER] Error encoding.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Infof("[HANDLER] Hello message sent.")
}
