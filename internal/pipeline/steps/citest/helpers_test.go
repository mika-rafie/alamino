package citest

import (
	"testing"

	"github.com/mika-rafie/alamino/internal/pipeline"
	"github.com/mika-rafie/alamino/internal/types"
)

func assertCIRestartsValidation(t *testing.T, outcome *pipeline.StepOutcome, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("CI repair returned error: %v", err)
	}
	if outcome == nil || outcome.RestartFrom != types.StepReview {
		t.Fatalf("CI repair outcome = %#v, want restart from review", outcome)
	}
}
