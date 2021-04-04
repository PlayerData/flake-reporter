package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"cloud.google.com/go/firestore"
	"playerdata.co.uk/flake-reporter/internal/models"
)

type TestSummaryHandler struct {
	Client *firestore.Client
	Ctx    context.Context
}

func (handler *TestSummaryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathComponents := strings.Split(r.URL.Path, "/")
	testName := pathComponents[len(pathComponents)-1]
	suiteName := pathComponents[len(pathComponents)-2]

	payload := models.TestSummary{
		Test:  testName,
		Suite: suiteName,
	}

	payload.PopulateFlakiness(handler.Client, handler.Ctx)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}
