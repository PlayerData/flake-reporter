package http

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/prometheus/client_golang/prometheus"
	"playerdata.co.uk/flake-reporter/internal/helpers"
	"playerdata.co.uk/flake-reporter/internal/models"
)

var testFlakinessGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "flake_reporter_test_flakiness", Help: "Flakiness for a test"},
	[]string{"project", "suite", "test", "branch"},
)

func setTestFlakiness(client *firestore.Client, ctx context.Context) {
	projects := client.Collection("projects").DocumentRefs(ctx)

	err := helpers.FirestoreForEachDocument(projects, func(project *firestore.DocumentRef) error {
		suites := project.Collection("suites").DocumentRefs(ctx)
		return helpers.FirestoreForEachDocument(suites, func(suite *firestore.DocumentRef) error {
			tests := suite.Collections(ctx)
			return helpers.FirestoreForEachCollection(tests, func(test *firestore.CollectionRef) error {
				branches := test.DocumentRefs(ctx)
				return helpers.FirestoreForEachDocument(branches, func(branch *firestore.DocumentRef) error {
					return client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
						summary, err := models.ReadBranchSummary(client, tx, project.ID, suite.ID, test.ID, branch.ID)
						if err != nil {
							return err
						}

						testFlakinessGauge.WithLabelValues(project.ID, suite.ID, test.ID, branch.ID).Set(float64(summary.Flakiness()))

						return nil
					})
				})
			})
		})
	})

	if err != nil {
		log.Printf("Failed to update Prometheus metrics: %s", err)
	}
}

func updateLoop(client *firestore.Client, ctx context.Context) {
	for {
		setTestFlakiness(client, ctx)

		time.Sleep(time.Duration(5) * time.Minute)
	}
}

func RegisterMetrics(client *firestore.Client, ctx context.Context) {
	prometheus.MustRegister(testFlakinessGauge)

	go updateLoop(client, ctx)
}
