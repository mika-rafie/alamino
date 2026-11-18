package main

import (
	"slices"
	"testing"
	"time"
)

func TestCIWorkflowRunsTestsOnAllSupportedDesktopPlatforms(t *testing.T) {
	job := ciTestJob(t)
	got := make(map[string]int)
	for _, row := range job.Strategy.Matrix.Include {
		got[row["os"]]++
	}

	want := map[string]int{
		"ubuntu-latest":  1,
		"macos-latest":   1,
		"windows-latest": 15,
	}
	if len(got) != len(want) {
		t.Fatalf("test matrix operating systems = %v, want %v", got, want)
	}
	for osName, count := range want {
		if got[osName] != count {
			t.Errorf("test matrix rows for %q = %d, want %d", osName, got[osName], count)
		}
	}
}

func TestCIWorkflowUsesRaceTestsOnUnixRunners(t *testing.T) {
	job := ciTestJob(t)
	commands := workflowCommandsMatching(job.Steps, func(step wfStep) bool {
		return exactRunnerOSCondition(step.If, "!=", "Windows")
	})

	var raceTests []workflowCommand
	for _, command := range commands {
		if command.name == "go" && len(command.args) > 0 && command.args[0] == "test" &&
			slices.Contains(command.args, "-race") && slices.Contains(command.args, "./...") {
			raceTests = append(raceTests, command)
		}
	}
	if len(raceTests) != 1 {
		t.Fatalf("Unix-only go test -race ./... commands = %d, want 1; normalized commands: %#v", len(raceTests), commands)
	}

	// The slowest race-instrumented packages sit at Go's default 10m
	// per-package limit on a cold macOS runner, so the limit must be raised
	// explicitly, and it must still fire before the job cap so a hang ends in
	// a goroutine dump instead of an evidence-free cancellation.
	goTimeout := goTestTimeout(t, raceTests[0])
	jobTimeout := time.Duration(job.TimeoutMinutes) * time.Minute
	if goTimeout <= 10*time.Minute {
		t.Fatalf("Unix go test -timeout is %s, want more than Go's 10m default", goTimeout)
	}
	if goTimeout >= jobTimeout {
		t.Fatalf("Unix go test -timeout is %s and the job cap is %s; the Go timeout must fire first", goTimeout, jobTimeout)
	}
}
