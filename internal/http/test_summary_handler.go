package http

import (
	"context"
	"encoding/json"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
	"playerdata.co.uk/flake-reporter/internal/models"
)

type TestSummaryHandler struct {
	Client *firestore.Client
	Ctx    context.Context
}

func (handler *TestSummaryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectName := vars["project"]
	suiteName := vars["suite"]
	testName := vars["test"]

	payload := models.TestSummary{
		Project: projectName,
		Suite:   suiteName,
		Test:    testName,
	}

	payload.PopulateFlakiness(handler.Client, handler.Ctx)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}
