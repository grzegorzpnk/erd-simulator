package api

import (
	log "10.254.188.33/matyspi5/erd/pkg/nmt/src/logger"
	"10.254.188.33/matyspi5/erd/pkg/nmt/src/pkg/model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func (h *apiHandler) deployApplication(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var mecApp model.MECApp
	_ = json.NewDecoder(r.Body).Decode(&mecApp)
	log.Infof("Client tries to deploy new mecApp ID: %v at MEC Host: %v \n", mecApp.Id, params["cluster"])
	if !(h.graphClient.CheckEnoughResources(mecApp, params["cluster"])) {
		err := fmt.Errorf("Mec App %v,cannot be instantiated beacuse of not enough resources", mecApp.Id)
		log.Errorf(err.Error())
		w.WriteHeader(http.StatusConflict)
	} else {
		h.graphClient.InstantiateApp(mecApp, params["cluster"])
		w.WriteHeader(http.StatusOK)
	}
}
