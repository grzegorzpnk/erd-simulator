package api

import (
	"net/http"
)

type apiHandler struct{}

func (h *apiHandler) SetClients() {}

//main functions
func (h *apiHandler) relocateApp(w http.ResponseWriter, r *http.Request) {

}

func (h *apiHandler) instantiateApps(w http.ResponseWriter, r *http.Request) {

}