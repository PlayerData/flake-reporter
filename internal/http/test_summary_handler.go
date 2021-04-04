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
	testName := vars["test"]
	suiteName := vars["suite"]

	payload := models.TestSummary{
		Test:  testName,
		Suite: suiteName,
	}

	payload.PopulateFlakiness(handler.Client, handler.Ctx)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}
