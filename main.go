package main

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
)

type TestStats struct {
	FailureCount int `firestore:"failure_count"`
	SuccessCount int `firestore:"success_count"`
}

func main() {
	ctx := context.Background()

	firestoreClient, err := firestore.NewClient(ctx, "test-project")
	if err != nil {
		log.Fatalf("firebase.NewClient err: %v", err)
	}
	defer firestoreClient.Close()

	http.Handle("/recv/junit", &JUnitHandler{client: firestoreClient, ctx: ctx})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
