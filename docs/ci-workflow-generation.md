# CI Workflow Generation

The `alamino ci-workflow` command auto-generates `.github/workflows/ci.yml` from your `.alamino.yaml` configuration, ensuring that GitHub checks align with the gate's canonical commands.

## Why This Matters

When a new repo is initialized with `alamino init`, the gate is ready to validate pushes locally — but GitHub has no CI workflow to register checks for. The gate's `ci` step sits idle (or times out) waiting for status checks that never appear.

This command generates those checks automatically from the source of truth: `.alamino.yaml`.

## Quick Start

After initializing your repo with `alamino init`:

```bash
alamino ci-workflow
```

This creates `.github/workflows/ci.yml` with steps that mirror your configured lint and test commands.

Commit and push the workflow file:

```bash
git add .github/workflows/ci.yml
git commit -m 'ci: register workflow'
git push alamino main
```

## Example

### Before

A freshly initialized Go repo with a root `go.mod`:

```bash
$ ls -la .github/workflows/
ls: .github/workflows/: No such file or directory

$ cat .alamino.yaml
commands:
  lint: "go vet ./... && gofmt -l ."
  test: "go test ./... -race"
```

Running `alamino ci-workflow`:

```bash
$ alamino ci-workflow
  ✓  CI workflow generated
      .github/workflows/ci.yml

  Commit and push this file to enable GitHub checks:
  git add .github/workflows/ci.yml && git commit -m 'ci: register workflow'
```

### After

GitHub Actions now registers real checks:

```yaml
name: CI

on:
  push:
    branches: ['main']
  pull_request:

permissions:
  contents: read

concurrency:
  group: ci-${{ github.ref }}
  cancel-in-progress: true

jobs:
  build-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
          check-latest: true

      - name: Lint
        run: |
          go vet ./... && gofmt -l .

      - name: Test
        run: |
          go test ./... -race
```

When you push through the alamino gate or open a PR, GitHub now has real checks to gate on. The gate's `ci` step completes successfully instead of timing out.

**Note:** The `push` branch is resolved from your repository's default branch (not hardcoded to `main`). The `pull_request` trigger runs on all PRs. Commands are placed in YAML block scalars so special characters and multi-line commands stay intact.

## Options

### Force Overwrite

If you need to regenerate the workflow (e.g., after changing `.alamino.yaml`), use `--force`:

```bash
alamino ci-workflow --force
```

## Notes

- See the [CLI reference](src/content/docs/reference/cli.md#alamino-ci-workflow) for toolchain setup behavior and customization.
- Commands are inserted verbatim from `.alamino.yaml`, so ensure they're portable to GitHub's Ubuntu runner.
- The workflow runs on push to the repository's default branch and on all `pull_request` events.
