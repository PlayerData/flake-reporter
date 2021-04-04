package models

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

type TestSummary struct {
	Suite         string
	Test          string
	BranchSummary []BranchResultSummary `json:"-"`
	Flakiness     float32
}

func testCollection(client *firestore.Client, suite string, test string) *firestore.CollectionRef {
	collectionPath := fmt.Sprintf("projects/test-project/suites/%s/%s", suite, test)
	return client.Collection(collectionPath)
}

func (summary *TestSummary) populateBranchResults(client *firestore.Client, ctx context.Context) error {
	collectionRef := testCollection(client, summary.Suite, summary.Test)
	branchResultSummaryRefs := collectionRef.Documents(ctx)

	branchDocuments, err := branchResultSummaryRefs.GetAll()
	if err != nil {
		return err
	}

	for _, branchDocument := range branchDocuments {
		var branchSummary BranchResultSummary
		err = branchDocument.DataTo(&branchSummary)
		if err != nil {
			return err
		}

		summary.BranchSummary = append(summary.BranchSummary, branchSummary)
	}

	return nil
}

func (summary *TestSummary) calculateFlakiness() float32 {
	var passCount int = 0
	var runCount int = 0

	for _, summary := range summary.BranchSummary {
		passCount += summary.passCount()
		runCount += summary.runCount()
	}

	passRate := float32(passCount) / float32(runCount)

	return 1.0 - passRate
}

func (summary *TestSummary) PopulateFlakiness(client *firestore.Client, ctx context.Context) error {
	err := summary.populateBranchResults(client, ctx)
	if err != nil {
		return err
	}

	summary.Flakiness = summary.calculateFlakiness()

	return nil
}
