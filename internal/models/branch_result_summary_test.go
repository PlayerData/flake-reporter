package models

import "testing"

func TestBranchResultSummaryFlakiness(t *testing.T) {
	summary := BranchResultSummary{
		Results: []int{1, 0, 1, 1},
	}

	flakiness := summary.flakiness()

	if flakiness != 0.25 {
		t.Fatalf("Flakiness reported as %v", flakiness)
	}
}
