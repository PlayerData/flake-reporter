package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"playerdata.co.uk/flake-reporter/internal/models"
	"playerdata.co.uk/flake-reporter/internal/test"
)

func TestTestSummaryHandler(t *testing.T) {
	ctx := context.Background()
	firestoreClient := test.NewFirestoreTestClient(ctx)
	test.ClearFirestore(t)

	mainSummary := models.BranchResultSummary{Results: []int{1, 1, 1, 0, 1}}
	featureBranchSummary := models.BranchResultSummary{Results: []int{0, 0, 0, 0, 1}}

	docRef := models.SummaryDocRef(firestoreClient, "test-suite", "example test", "main")
	docRef.Set(ctx, mainSummary)

	docRef = models.SummaryDocRef(firestoreClient, "test-suite", "example test", "feature-branch")
	_, err := docRef.Set(ctx, featureBranchSummary)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/summary/projects/test-project/suites/test-suite/example test", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	handler := http.Handler(&TestSummaryHandler{Client: firestoreClient, Ctx: ctx})
	handler.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v", status)
		t.Fatalf("reponse body: got %s", res.Body.String())
	}

	var responsePayload models.TestSummary
	json.Unmarshal(res.Body.Bytes(), &responsePayload)

	expectedPayload := models.TestSummary{
		Suite:     "test-suite",
		Test:      "example test",
		Flakiness: 0.5,
	}

	if !reflect.DeepEqual(responsePayload, expectedPayload) {
		t.Fatalf("Got incorrect response: %v", responsePayload)
	}
}
