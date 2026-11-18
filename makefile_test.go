package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// The build stamps version metadata only. No analytics endpoint or site ID is
// embedded, whatever the environment or a local .env file says.
func TestMakeBuildStampsVersionMetadataOnly(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("make build tests are POSIX-oriented")
	}
	makePath, err := exec.LookPath("make")
	if err != nil {
		t.Skip("make not available")
	}

	data, err := os.ReadFile("Makefile")
	if err != nil {
		t.Fatal(err)
	}
	workDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(workDir, "Makefile"), data, 0o644); err != nil {
		t.Fatal(err)
	}
	dotenv := "VERSION=from-dotenv\nALAMINO_UMAMI_HOST=https://dotenv.example\nALAMINO_UMAMI_WEBSITE_ID=website-from-dotenv\n"
	if err := os.WriteFile(filepath.Join(workDir, ".env"), []byte(dotenv), 0o644); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, makePath, "-n", "build")
	cmd.Dir = workDir
	cmd.Env = append(os.Environ(), "UMAMI_HOST=https://env.example", "UMAMI_WEBSITE_ID=website-from-env")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("make -n build failed: %v\n%s", err, out)
	}
	output := string(out)

	for _, want := range []string{"buildinfo.Version=", "buildinfo.Commit=", "buildinfo.Date="} {
		if !strings.Contains(output, want) {
			t.Errorf("make build must stamp %s, got:\n%s", want, output)
		}
	}
	for _, forbidden := range []string{"Telemetry", "dotenv.example", "env.example", "website-from-", "buildinfo.Version=from-dotenv"} {
		if strings.Contains(output, forbidden) {
			t.Errorf("make build must not embed %q, got:\n%s", forbidden, output)
		}
	}
}
