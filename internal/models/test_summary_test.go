package models

import (
	"context"
	"reflect"
	"testing"

	"playerdata.co.uk/flake-reporter/internal/test"
)

func TestSummaryPopulateBranchResults(t *testing.T) {
	ctx := context.Background()
	firestoreClient := test.NewFirestoreTestClient(ctx)
	test.ClearFirestore(t)

	mainSummary := BranchResultSummary{Results: []int{1, 1, 1, 0, 1}}
	featureBranchSummary := BranchResultSummary{Results: []int{1, 1, 1, 0, 1}}

	docRef := SummaryDocRef(firestoreClient, "test-suite", "example test", "main")
	docRef.Set(ctx, mainSummary)

	docRef = SummaryDocRef(firestoreClient, "test-suite", "example test", "feature-branch")
	docRef.Set(ctx, featureBranchSummary)

	summary := TestSummary{Suite: "test-suite", Test: "example test"}

	err := summary.populateBranchResults(firestoreClient, ctx)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(summary.BranchSummary, []BranchResultSummary{mainSummary, featureBranchSummary}) {
		t.Fatalf("branch results wrong: %v", summary.BranchSummary)
	}
}

func TestSummaryFlakiness(t *testing.T) {
	summary := TestSummary{
		BranchSummary: []BranchResultSummary{
			{Results: []int{1, 0, 1, 1}},
			{Results: []int{1, 1, 1, 1}},
			{Results: []int{1, 1, 1, 1}},
			{Results: []int{1, 1, 1, 1}},
		},
	}

	flakiness := summary.calculateFlakiness()

	if flakiness != 0.0625 {
		t.Fatalf("Flakiness reported as %v", flakiness)
	}
}
