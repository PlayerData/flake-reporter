package http

import (
	"context"
	"math"
	"testing"

	prometheusModel "github.com/prometheus/client_model/go"

	"playerdata.co.uk/flake-reporter/internal/models"
	"playerdata.co.uk/flake-reporter/internal/test"
)

func TestPrometheusHandler(t *testing.T) {
	ctx := context.Background()
	firestoreClient := test.NewFirestoreTestClient(ctx)
	test.ClearFirestore(t)

	mainSummary := models.BranchResultSummary{Results: []int{1, 1, 1, 0, 1}}
	featureBranchSummary := models.BranchResultSummary{Results: []int{0, 0, 0, 0, 1}}

	docRef := models.SummaryDocRef(firestoreClient, "test-project", "test-suite", "example test", "main")
	docRef.Set(ctx, mainSummary)

	docRef = models.SummaryDocRef(firestoreClient, "test-project", "test-suite", "example test", "feature-branch")
	docRef.Set(ctx, featureBranchSummary)

	setTestFlakiness(firestoreClient, ctx)

	result, err := testFlakinessGauge.GetMetricWithLabelValues("test-project", "test-suite", "example test", "main")
	if err != nil {
		t.Fatal(err)
	}

	metric := &prometheusModel.Metric{}

	err = result.Write(metric)
	if err != nil {
		t.Fatal(err)
	}

	value := metric.GetGauge().GetValue()
	delta := 0.2 - value

	if math.Abs(delta) > 0.001 {
		t.Fatalf("Value was %f", value)
	}
}
