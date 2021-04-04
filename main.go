package main

import (
	"context"
	"log"
	nethttp "net/http"

	"cloud.google.com/go/firestore"
	"playerdata.co.uk/flake-reporter/internal/http"
)

func main() {
	ctx := context.Background()

	firestoreClient, err := firestore.NewClient(ctx, "test-project")
	if err != nil {
		log.Fatalf("firebase.NewClient err: %v", err)
	}
	defer firestoreClient.Close()

	nethttp.Handle("/recv/junit", &http.JUnitHandler{Client: firestoreClient, Ctx: ctx})
	nethttp.Handle("/summary", &http.TestSummaryHandler{Client: firestoreClient, Ctx: ctx})

	log.Fatal(nethttp.ListenAndServe(":8080", nil))
}
