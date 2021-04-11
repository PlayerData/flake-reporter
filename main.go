package main

import (
	"context"
	"log"
	nethttp "net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
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

	r := mux.NewRouter()
	r.Handle("/recv/junit", &http.JUnitHandler{Client: firestoreClient, Ctx: ctx})
	r.Handle("/summary/{project}/{suite}/{test}", &http.TestSummaryHandler{Client: firestoreClient, Ctx: ctx})

	nethttp.Handle("/", r)

	log.Fatal(nethttp.ListenAndServe(":9090", nil))
}
