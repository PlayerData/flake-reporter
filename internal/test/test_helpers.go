package test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
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
	emulatorHost := os.Getenv("FIRESTORE_EMULATOR_HOST")
	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://%s/emulator/v1/projects/test-project/databases/(default)/documents", emulatorHost), nil)
	if err != nil {
		t.Fatal(err)

	}

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
}
