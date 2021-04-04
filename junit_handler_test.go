package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"cloud.google.com/go/firestore"
)

func newFirestoreTestClient(ctx context.Context) *firestore.Client {
	client, err := firestore.NewClient(ctx, "test-project")
	if err != nil {
		log.Fatalf("firebase.NewClient err: %v", err)
	}

	return client
}

func clearFirestore(t *testing.T) {
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

func prepareFileUpload(t *testing.T, writer *multipart.Writer, filename string) {
	path := "./fixtures/" + filename
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}

	defer file.Close()
	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Fatal(err)
	}
	io.Copy(part, file)
	writer.Close()
}

func uploadTestResult(t *testing.T, firestoreClient *firestore.Client, ctx context.Context, fixtureName string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("project", "test-project")
	writer.WriteField("branch", "master")
	prepareFileUpload(t, writer, fixtureName)

	req, err := http.NewRequest("POST", "/recv/junit", body)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())

	res := httptest.NewRecorder()

	handler := http.Handler(&JUnitHandler{client: firestoreClient, ctx: ctx})
	handler.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v", status)
		t.Fatalf("reponse body: got %s", res.Body.String())
	}
}

func TestUpdateTestSummary(t *testing.T) {
	ctx := context.Background()
	firestoreClient := newFirestoreTestClient(ctx)
	clearFirestore(t)

	uploadTestResult(t, firestoreClient, ctx, "junit-success.xml")
	uploadTestResult(t, firestoreClient, ctx, "junit-success.xml")
	uploadTestResult(t, firestoreClient, ctx, "junit-failure.xml")
	uploadTestResult(t, firestoreClient, ctx, "junit-success.xml")

	err := firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		summary, err := readTestSummary(firestoreClient, tx, "Session Tags", "Session Tags can create a session tag", "master")
		if err != nil {
			t.Error(err)
			return err
		}

		if !reflect.DeepEqual(summary.Results, []int{1, 1, 0, 1}) {
			return fmt.Errorf("summary results wrong: %v", summary.Results)
		}

		return nil
	})

	if err != nil {
		t.Fatal(err)
	}
}
