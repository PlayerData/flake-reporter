package models

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/joshdk/go-junit"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BranchResultSummary struct {
	Results []int `firestore:"results"`
}

func (summary *BranchResultSummary) passCount() int {
	var passes int = 0

	for _, result := range summary.Results {
		passes += result
	}

	return passes
}

func (summary *BranchResultSummary) runCount() int {
	return len(summary.Results)
}

func (summary *BranchResultSummary) flakiness() float32 {
	passRate := float32(summary.passCount()) / float32(summary.runCount())

	return 1.0 - passRate
}

func SummaryDocRef(client *firestore.Client, suite string, test string, branch string) *firestore.DocumentRef {
	return testCollection(client, suite, test).Doc(branch)
}

func ReadBranchSummary(client *firestore.Client, tx *firestore.Transaction, suite string, test string, branch string) (BranchResultSummary, error) {
	var summary BranchResultSummary
	docRef := SummaryDocRef(client, suite, test, branch)

	doc, err := tx.Get(docRef)
	if err != nil {
		return BranchResultSummary{}, err
	}

	err = doc.DataTo(&summary)
	if err != nil {
		return BranchResultSummary{}, err
	}

	return summary, nil
}

func resultInt(testResult junit.Test) int {
	if testResult.Status == junit.StatusPassed {
		return 1
	} else {
		return 0
	}
}

func UpdateBranchSummary(client *firestore.Client, ctx context.Context, project string, suite string, branch string, testResult junit.Test) error {
	err := client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		docRef := SummaryDocRef(client, suite, testResult.Name, branch)

		summary, err := ReadBranchSummary(client, tx, suite, testResult.Name, branch)
		if err != nil && status.Code(err) != codes.NotFound {
			return err
		}

		summary.Results = append(summary.Results, resultInt(testResult))

		return tx.Set(docRef, summary)
	})

	return err
}
