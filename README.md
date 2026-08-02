# 🏔️ Atlas Platform

> A personal project to build and operate a production-grade platform from scratch — a Go REST API packaged with Docker, validated through CI/CD, deployed with Kubernetes and Helm, managed through a GitOps workflow, and backed by observability.

---

## 📋 Table of Contents

- [About](#about)
- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Folder Structure](#folder-structure)
- [Project Roadmap](#project-roadmap)
- [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Run Microservice Locally](#run-microservice-locally)
- [Local Setup](#local-setup)
- [API Reference](#api-reference)
- [Environment Variables](#environment-variables)

---

## About

**Atlas Platform** is a hands-on platform engineering project built end-to-end from the ground up. The goal is to take a Go REST API from a developer's machine all the way to a production-ready, observable, and automatically deployed system — one milestone at a time.

Check the [Project Roadmap](#project-roadmap) to see where things currently stand.

---

## Architecture

![Atlas Platform Architecture](./docs/architecture.png)

---

## Tech Stack

| Layer | Technology                                              |
|---|---------------------------------------------------------|
| **Language** | Go                                                      |
| **Router** | [Chi](https://github.com/go-chi/chi)                    |
| **Database** | PostgreSQL                                              |
| **Migrations** | [Goose](https://github.com/pressly/goose)               |
| **Containerisation** | Docker                                                  |
| **CI Pipeline** | GitHub Actions                                          |
| **Orchestration** | Kubernetes + Helm                                       |
| **GitOps** | ArgoCD                                                  |
| **Observability** | Grafana Ecosystem (Prometheus + Loki + Tempo + Grafana) |

---

## Folder Structure

| Path | Purpose |
|---|---|
| `cmd/api` | API entrypoint, HTTP server setup, route handlers, JSON helpers, error responses, and API tests |
| `cmd/migrate/migrations` | Goose database migration files |
| `internal/db` | PostgreSQL connection setup using `pgxpool` |
| `internal/env` | Environment variable loading and fallback helpers |
| `internal/store` | Data access layer and user storage operations |
| `K8s/manifest` | Raw Kubernetes manifests for ConfigMap, Secret example, PostgreSQL, migration job, API deployment, service, and ingress |
| `helm/atlas-platform-chart` | Helm chart used to deploy the platform to Kubernetes |
| `gitops/argocd` | Argo CD setup notes for GitOps workflow |
| `docs` | Project documentation assets, including architecture diagrams |
| `.github/workflows` | GitHub Actions CI/CD pipeline |

---

## Project Roadmap

| # | Milestone | Status |
|---|---|---|
| 1 | Create a simple REST API Webserver | ✅ Done |
| 2 | Containerise REST API | ✅ Done |
| 3 | Setup one-click local development setup | ✅ Done  |
| 4 | Setup a CI pipeline | ✅ Done  |
| 5 | Deploy in Kubernetes | ✅ Done |
| 6 | Deploy using Helm Charts | ✅ Done |
| 7 | Setup GitOps with ArgoCD | ✅ Done |
| 8 | Setup observability stack | ⬜ Upcoming |
| 9 | Configure dashboards & alerts | ⬜ Upcoming |

---

## Getting Started

### Prerequisites - Tools + Docs

- [Go 1.25+](https://go.dev/dl/)
- [PostgreSQL](https://www.postgresql.org/download/)
- [Goose](https://github.com/pressly/goose) — for running migrations
- Docker - Install Locally (https://docker.com)
- Kubernetes - Local cluster for dev : Anyone of the following work
  - Minikube (https://minikube.sigs.k8s.io/docs/)
  - Kind (https://kind.sigs.k8s.io/)
- AWS EKS - Cloud K8s cluster
  - AWS CLI (https://aws.amazon.com/cli/)
  - eksctl (https://docs.aws.amazon.com/eks/latest/userguide/getting-started-eksctl.html)
- Github Actions (https://docs.github.com/en/actions)
- ArgoCD (https://argo-cd.readthedocs.io/)

### Run Microservice Locally

Clone the repository:

```bash
git clone https://github.com/<your-username>/atlas-platform.git
cd atlas-platform
```

Set up environment variables:

```bash
cp .env.example .env
# Edit .env with your local DB credentials and config
```

Use your own local values in `.env`. The example file is only a template, and real secrets should stay out of git.

Install dependencies:

```bash
go mod tidy
```

Run database migrations:

```bash
goose -dir cmd/migrate/migrations postgres "$DATABASE_URL" up
```

Start the API server:

```bash
go run ./cmd/api
```

The server will be available at `http://localhost:3030`.

Health check:

```bash
curl http://localhost:3030/health
# {"data":{"status":"ok","env":"dev","version":"v1"}}
```

Run tests:

```bash
go test ./...
```

---

## Local Setup

### Docker Compose

Create your local `.env` from `.env.example` before starting the Compose stack.

Use the Taskfile to build the API image, start PostgreSQL, run migrations, and start the API container:

```bash
task up
```

Clean up the local Compose setup, including containers, networks, and volumes:

```bash
task clean
```

### Kubernetes Manifests

Make sure your local Kubernetes cluster is running and your `kubectl` context points to it. Create your own Kubernetes Secret manifest from `K8s/manifest/secrets.example.yaml`, then deploy the platform using the raw Kubernetes manifests:

```bash
task deploy
```

### Helm

Create your own Kubernetes Secret values from the example secret manifest before installing the chart. Then deploy the platform using Helm:

```bash
helm upgrade --install atlas-platform ./helm/atlas-platform-chart
```

---

## API Reference

### Health

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/health` | Returns server health status |

**Response:**
```json
{
  "data": {
    "status": "ok",
    "env": "dev",
    "version": "v1"
  }
}
```

### Users

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/v1/users/` | Creates a new user |
| `GET` | `/v1/users/` | Returns all users |
| `GET` | `/v1/users/{id}/` | Returns a user by ID |
| `PATCH` | `/v1/users/{id}/` | Updates a user by ID |
| `DELETE` | `/v1/users/{id}/` | Deletes a user by ID |

**Create user request:**
```json
{
  "name": "Ankur",
  "age": 24,
  "email": "ankur@example.com"
}
```

**Update user request:**
```json
{
  "name": "Updated Name",
  "age": 25,
  "email": "updated@example.com"
}
```

All successful responses are wrapped in a `data` object.

---

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `DATABASE_URL` | — | PostgreSQL connection string used by the API |
| `PORT` | `:3030` | Address and port the API server listens on |
| `ENV` | `dev` | Runtime environment |
| `VERSION` | `v1` | API version returned by the health check |
| `EXTERNAL_URL` | `localhost:3030` | External API URL used by application config |
| `GOOSE_DRIVER` | `postgres` | Goose database driver |
| `GOOSE_DBSTRING` | — | PostgreSQL connection string used by Goose migrations |
| `GOOSE_MIGRATION_DIR` | `/app/cmd/migrate/migrations` | Migration directory used by Goose in containers and Kubernetes |
| `POSTGRES_USER` | — | PostgreSQL container username |
| `POSTGRES_PASSWORD` | — | PostgreSQL container password |
| `POSTGRES_DB` | — | PostgreSQL container database name |

---

> **Note:** This is a personal project, built incrementally as a learning exercise in platform and reliability engineering. Every milestone adds a new layer to the system — follow along as it grows.
