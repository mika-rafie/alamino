---
name: release-workflow
description: Use when changing the release workflow, release artifact names or checksums, or the self-update channel manifest.
user-invocable: false
metadata:
  internal: true
---

**Release artifacts (`.github/workflows/release.yml`)**

- One ubuntu job (`build-and-upload`) cross-compiles all six targets (`darwin`, `linux`, `windows` x `amd64`, `arm64`) with `CGO_ENABLED=0`, and stamps only `buildinfo.Version`, `Commit`, and `Date`.
- macOS binaries are not Developer ID signed or notarized, and the workflow needs no macOS runner, protected environment, repository secret, or Actions variable. The Go linker ad-hoc signs `darwin/arm64`, which is what Apple silicon requires to run it; `alamino update` clears the quarantine attribute after replacing the binary (`internal/update/quarantine_darwin.go`).
- An ad-hoc signature has a content-based designated requirement, so macOS permission grants (TCC) do not carry over from one release to the next.
- The archive names `alamino-<tag>-<goos>-<goarch>.tar.gz` (`.zip` on Windows) and `checksums.txt` (`sha256sum alamino-*`) are the installer and `alamino update` contract; `finalize` publishes the draft as a prerelease only after every asset job succeeds.
- The contract is pinned by the root `TestReleaseWorkflow*` static tests in `workflow_release_test.go`.

**Self-update channel manifest (`internal/update`)**

- `alamino update` reads version metadata exclusively from `channels.json` on the GitHub release-asset CDN (`releases/download/channels/channels.json`), not `api.github.com`; a token is never required. Publisher: `cmd/publish-channels`, invoked from `.github/workflows/publish-channels.yml` (reusable `workflow_call`, plus `workflow_dispatch` / `on: release` backstops). `release.yml` calls it after `finalize` because GitHub does not cascade `GITHUB_TOKEN` `release` events, so release-please (pre)releases would otherwise leave the channel stale. Regressions: `internal/update/channels_test.go`, `workflow_publish_channels_test.go`, `TestReleaseWorkflowCallsPublishChannelsAfterFinalize`.
