package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"cloud.google.com/go/firestore"

	"playerdata.co.uk/flake-reporter/internal/models"
	"playerdata.co.uk/flake-reporter/internal/test"
)

func prepareFileUpload(t *testing.T, writer *multipart.Writer, filename string) {
	path := "../../fixtures/" + filename
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

	handler := http.Handler(&JUnitHandler{Client: firestoreClient, Ctx: ctx})
	handler.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v", status)
		t.Fatalf("reponse body: got %s", res.Body.String())
	}
}

func TestUpdateTestSummary(t *testing.T) {
	ctx := context.Background()
	firestoreClient := test.NewFirestoreTestClient(ctx)
	test.ClearFirestore(t)

	uploadTestResult(t, firestoreClient, ctx, "junit-success.xml")
	uploadTestResult(t, firestoreClient, ctx, "junit-success.xml")
	uploadTestResult(t, firestoreClient, ctx, "junit-failure.xml")
	uploadTestResult(t, firestoreClient, ctx, "junit-success.xml")

	err := firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		summary, err := models.ReadBranchSummary(firestoreClient, tx, "test-project", "Session Tags", "Session Tags can create a session tag", "master")
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
