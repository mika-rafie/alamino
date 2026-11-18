package steps

import (
	"context"
	"strings"
	"testing"

	"github.com/mika-rafie/alamino/internal/config"
	"github.com/mika-rafie/alamino/internal/db"
	"github.com/mika-rafie/alamino/internal/pipeline"
	"github.com/mika-rafie/alamino/internal/scm"
)

func TestDetectProviderForStep_CursorOriginGitRemote(t *testing.T) {
	sctx := &pipeline.StepContext{Ctx: context.Background()}
	got := detectProviderForStep(sctx, "https://origin.cursor.com/owner/repo.git")
	if got != scm.ProviderOrigin {
		t.Fatalf("detectProviderForStep() = %q, want %q", got, scm.ProviderOrigin)
	}
}

func TestDetectProviderForStep_CursorOriginWebPR(t *testing.T) {
	sctx := &pipeline.StepContext{Ctx: context.Background()}
	got := detectProviderForStep(sctx, "https://cursor.com/codebase/owner/repo/pull/6")
	if got != scm.ProviderOrigin {
		t.Fatalf("detectProviderForStep(web PR) = %q, want %q", got, scm.ProviderOrigin)
	}
}

func TestBuildHost_Origin(t *testing.T) {
	sctx := &pipeline.StepContext{
		Ctx: context.Background(),
		Run: &db.Run{Branch: "feat/x", HeadSHA: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		Repo: &db.Repo{
			UpstreamURL:   "https://origin.cursor.com/owner/repo.git",
			DefaultBranch: "master",
		},
		Config: &config.Config{},
	}

	host, reason := buildHost(sctx, scm.ProviderOrigin)
	if host == nil || reason != "" {
		t.Fatalf("buildHost() = (%v, %q), want Origin host", host, reason)
	}
	if host.Provider() != scm.ProviderOrigin {
		t.Fatalf("Provider() = %q, want %q", host.Provider(), scm.ProviderOrigin)
	}
}

func TestBuildHost_OriginResolvesSlugFromPRURL(t *testing.T) {
	prURL := "https://cursor.com/codebase/owner/repo/pull/6"
	sctx := &pipeline.StepContext{
		Ctx:    context.Background(),
		Run:    &db.Run{PRURL: &prURL},
		Repo:   &db.Repo{},
		Config: &config.Config{},
	}

	host, reason := buildHost(sctx, scm.ProviderOrigin)
	if host == nil || reason != "" {
		t.Fatalf("buildHost() = (%v, %q), want Origin host from PR URL", host, reason)
	}
}

func TestBuildHost_OriginRejectsForkRouting(t *testing.T) {
	sctx := &pipeline.StepContext{
		Ctx:  context.Background(),
		Run:  &db.Run{},
		Repo: &db.Repo{UpstreamURL: "https://origin.cursor.com/owner/repo.git", ForkURL: "https://origin.cursor.com/alice/repo.git"},
	}

	host, reason := buildHost(sctx, scm.ProviderOrigin)
	if host != nil || !strings.Contains(reason, "fork") {
		t.Fatalf("buildHost() = (%v, %q), want explicit fork rejection", host, reason)
	}
}
