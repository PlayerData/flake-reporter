package test

import (
	"context"
	"log"
	"net/http"
	"testing"

	"cloud.google.com/go/firestore"
)

func NewFirestoreTestClient(ctx context.Context) *firestore.Client {
	client, err := firestore.NewClient(ctx, "test-project")
	if err != nil {
		log.Fatalf("firebase.NewClient err: %v", err)
	}

	return client
}

func ClearFirestore(t *testing.T) {
	req, err := http.NewRequest("DELETE", "http://localhost:8080/emulator/v1/projects/test-project/databases/(default)/documents", nil)
	if err != nil {
		t.Fatal(err)

	}

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
}
