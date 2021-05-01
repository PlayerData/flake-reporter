package main

import (
	"context"
	"log"
	nethttp "net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"playerdata.co.uk/flake-reporter/internal/http"
)

func main() {
	ctx := context.Background()

	firestoreProject := os.Getenv("FIRESTORE_PROJECT_ID")
	if firestoreProject == "" {
		log.Fatalf("FIRESTORE_PROJECT_ID unset")
	}
	firestoreClient, err := firestore.NewClient(ctx, firestoreProject)
	if err != nil {
		log.Fatalf("firebase.NewClient err: %v", err)
	}
	defer firestoreClient.Close()

	http.RegisterMetrics(firestoreClient, ctx)

	r := mux.NewRouter()
	r.Handle("/metrics", promhttp.Handler())
	r.Handle("/recv/junit", &http.JUnitHandler{Client: firestoreClient, Ctx: ctx})
	r.Handle("/summary/{project}/{suite}/{test}", &http.TestSummaryHandler{Client: firestoreClient, Ctx: ctx})

	nethttp.Handle("/", r)

	log.Fatal(nethttp.ListenAndServe(":9090", nil))
}
