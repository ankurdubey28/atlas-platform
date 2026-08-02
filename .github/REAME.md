# GitHub Actions CI

This directory contains the GitHub Actions workflow for Atlas Platform.

## Workflow File

| File | Purpose |
|---|---|
| `workflows/ci.yml` | Runs CI checks, builds Docker images, pushes images to GHCR, and updates Helm image tags |

## When It Runs

The workflow runs on:

- Pushes to the `main` branch when API, migration, Docker, or Go dependency files change
- Manual runs from the GitHub Actions UI using `workflow_dispatch`

Tracked paths:

```yaml
cmd/**
internal/**
Dockerfile
go.mod
go.sum
```

## How The Workflow Works

### 1. Detect Changes

The `changes` job checks which parts of the repository changed.

It produces two outputs:

| Output | Meaning |
|---|---|
| `api` | API-related code, internal packages, Go dependencies, or Dockerfile changed |
| `migration` | Migration files or Dockerfile changed |

These outputs are used by later jobs to decide which Docker images need to be rebuilt.

### 2. Build And Test

The `build-test` job validates the Go project.

It does the following:

- Checks out the repository
- Sets up Go using the version from `go.mod`
- Downloads dependencies
- Checks formatting with `gofmt`
- Runs linting with `golangci-lint`
- Builds all packages
- Runs all tests

If this job fails, Docker images are not pushed.

### 3. Build And Push Images

The `push` job builds and pushes Docker images to GitHub Container Registry.

It pushes:

| Image | Docker target | Tag |
|---|---|---|
| `ghcr.io/<owner>/atlas-api` | `api` | GitHub run ID |
| `ghcr.io/<owner>/atlas-migrate` | `migrate` | GitHub run ID |

The API image is pushed when API-related files change.

The migration image is pushed when migration-related files change.

On a manual workflow run, both images can be built and pushed.

### 4. Update Helm Values

The `update-helm` job updates image tags in:

```text
helm/atlas-platform-chart/values.yaml
```

It updates:

| Helm value | Source |
|---|---|
| `.images.api.tag` | Current GitHub run ID |
| `.images.migrate.tag` | Current GitHub run ID |

After updating the file, the workflow commits and pushes the Helm values change back to the repository.

The commit message includes `[skip ci]` so the tag update commit does not start another CI run.

## Required Permissions And Secrets

The workflow uses:

| Item | Purpose |
|---|---|
| `contents: write` | Allows the workflow to commit updated Helm values |
| `packages: write` | Allows the workflow to push Docker images to GHCR |
| `GITHUB_TOKEN` | Used to authenticate with GHCR |
| `CI_BOT_NAME` | Git author name for Helm values update commits |
| `CI_BOT_EMAIL` | Git author email for Helm values update commits |

## Deployment Flow

The workflow does not directly deploy to Kubernetes.

Instead, it:

1. Validates the application
2. Builds and pushes versioned container images
3. Updates Helm chart image tags
4. Pushes the Helm values change

Argo CD can then detect the Helm chart change and sync the updated version into the Kubernetes cluster.
