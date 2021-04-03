package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func newFirestoreTestClient(ctx context.Context) *firestore.Client {
	client, err := firestore.NewClient(ctx, "test-project")
	if err != nil {
		log.Fatalf("firebase.NewClient err: %v", err)
	}

	return client
}

func clearFirestore(ctx context.Context, client *firestore.Client) error {
	collection := client.Collection("projects")

	for {
		// Get a batch of documents
		iter := collection.Documents(ctx)
		numDeleted := 0

		// Iterate through the documents, adding
		// a delete operation for each one to a
		// WriteBatch.
		batch := client.Batch()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}

			batch.Delete(doc.Ref)
			numDeleted++
		}

		// If there are no documents to delete,
		// the process is over.
		if numDeleted == 0 {
			return nil
		}

		_, err := batch.Commit(ctx)
		if err != nil {
			return err
		}
	}
}

func prepareFileUpload(t *testing.T, writer *multipart.Writer, filename string) {
	path := "./fixtures/" + filename
	file, err := os.Open(path)
	if err != nil {
		t.Error(err)
	}

	defer file.Close()
	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()
}

func TestReceiveJUnitNewSuccessReport(t *testing.T) {
	ctx := context.Background()
	firestoreClient := newFirestoreTestClient(ctx)

	clearFirestore(ctx, firestoreClient)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("project", "test-project")
	writer.WriteField("branch", "master")
	writer.WriteField("build-number", "1")
	prepareFileUpload(t, writer, "junit-success.xml")

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
		t.Errorf("reponse body: got %s", res.Body.String())
		return
	}

	collection := firestoreClient.Collection("projects/test-project/branches/master/suites/Session Tags/tests")

	resultDoc, err := collection.Doc("Session Tags can create a session tag").Get(ctx)
	if err != nil {
		t.Error(err)
		return
	}

	var result TestStats
	err = resultDoc.DataTo(&result)
	if err != nil {
		t.Error(err)
		return
	}

	if result.SuccessCount != 1 {
		t.Errorf("incorrect success count %v", result.SuccessCount)
		return
	}
}
