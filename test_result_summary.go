package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/joshdk/go-junit"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TestResultSummary struct {
	Results []int `firestore:"results"`
}

func summaryDocRef(client *firestore.Client, suite string, test string, branch string) *firestore.DocumentRef {
	docPath := fmt.Sprintf("projects/test-project/suites/%s/%s/%s", suite, test, branch)
	return client.Doc(docPath)
}

func readTestSummary(client *firestore.Client, tx *firestore.Transaction, suite string, test string, branch string) (TestResultSummary, error) {
	var summary TestResultSummary
	docRef := summaryDocRef(client, suite, test, branch)

	doc, err := tx.Get(docRef)
	if err != nil {
		return TestResultSummary{}, err
	}

	err = doc.DataTo(&summary)
	if err != nil {
		return TestResultSummary{}, err
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

func updateTestSummary(client *firestore.Client, ctx context.Context, project string, suite string, branch string, testResult junit.Test) error {
	err := client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		docRef := summaryDocRef(client, suite, testResult.Name, branch)

		summary, err := readTestSummary(client, tx, suite, testResult.Name, branch)
		if err != nil && status.Code(err) != codes.NotFound {
			return err
		}

		summary.Results = append(summary.Results, resultInt(testResult))

		return tx.Set(docRef, summary)
	})

	return err
}
