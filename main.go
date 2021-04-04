package main

import (
	"context"
	"log"
	nethttp "net/http"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
	"playerdata.co.uk/flake-reporter/internal/http"
)

func main() {
	ctx := context.Background()

	firestoreClient, err := firestore.NewClient(ctx, "test-project")
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
