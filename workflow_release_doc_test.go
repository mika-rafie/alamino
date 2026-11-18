package main

import (
	"os"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// wfDoc is a minimal typed view of a GitHub Actions workflow, shared by the
// static workflow tests in this package. Flexible fields (runs-on, environment,
// needs) are decoded as any and normalized by the helpers below.
type wfDoc struct {
	Jobs map[string]*wfJob `yaml:"jobs"`
}

type wfJob struct {
	name           string
	RunsOn         any    `yaml:"runs-on"`
	Environment    any    `yaml:"environment"`
	Needs          any    `yaml:"needs"`
	If             string `yaml:"if"`
	TimeoutMinutes int    `yaml:"timeout-minutes"`
	Permissions    struct {
		Contents string `yaml:"contents"`
	} `yaml:"permissions"`
	Uses     string     `yaml:"uses"`
	Strategy wfStrategy `yaml:"strategy"`
	Steps    []wfStep   `yaml:"steps"`
}

type wfStrategy struct {
	Matrix struct {
		Include []map[string]string `yaml:"include"`
	} `yaml:"matrix"`
}

type wfStep struct {
	Name  string            `yaml:"name"`
	Uses  string            `yaml:"uses"`
	If    string            `yaml:"if"`
	Shell string            `yaml:"shell"`
	Env   map[string]string `yaml:"env"`
	Run   string            `yaml:"run"`
	With  map[string]string `yaml:"with"`
}

func loadReleaseWorkflowDoc(t *testing.T) *wfDoc {
	t.Helper()
	raw, err := os.ReadFile(".github/workflows/release.yml")
	if err != nil {
		t.Fatalf("read release workflow: %v", err)
	}
	var wf wfDoc
	if err := yaml.Unmarshal(raw, &wf); err != nil {
		t.Fatalf("parse release workflow: %v", err)
	}
	for name, job := range wf.Jobs {
		job.name = name
	}
	return &wf
}

func wfScalar(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case map[string]any:
		if n, ok := t["name"].(string); ok {
			return n
		}
	}
	return ""
}

func wfList(v any) []string {
	switch t := v.(type) {
	case string:
		return []string{t}
	case []any:
		out := make([]string, 0, len(t))
		for _, e := range t {
			if s, ok := e.(string); ok {
				out = append(out, s)
			}
		}
		return out
	}
	return nil
}

func (j *wfJob) runsOn() string  { return wfScalar(j.RunsOn) }
func (j *wfJob) environ() string { return wfScalar(j.Environment) }
func (j *wfJob) needs() []string { return wfList(j.Needs) }

func (j *wfJob) allRun() string {
	var b strings.Builder
	for _, s := range j.Steps {
		b.WriteString(s.Run)
		b.WriteString("\n")
	}
	return b.String()
}
