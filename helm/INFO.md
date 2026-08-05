# Atlas Platform Helm Chart

This directory contains the Helm deployment workflow for Atlas Platform. The chart deploys the Go API, PostgreSQL, the database migration job, shared configuration, secrets, and ingress resources into a Kubernetes cluster.

## Directory Layout

| Path | Purpose |
|---|---|
| `Taskfile.yaml` | Local helper tasks for creating a Minikube cluster before installing the chart |
| `atlas-platform-chart/Chart.yaml` | Helm chart metadata |
| `atlas-platform-chart/values.yaml` | Configurable image values used by the templates |
| `atlas-platform-chart/templates/cm.yaml` | Runtime ConfigMap for API and migration settings |
| `atlas-platform-chart/templates/secrets.example.yaml` | Example secret template for creating your own credentials |
| `atlas-platform-chart/templates/statefulset.yaml` | PostgreSQL headless Service and StatefulSet |
| `atlas-platform-chart/templates/job.yaml` | Goose migration Job |
| `atlas-platform-chart/templates/deployment.yaml` | Atlas API Service and Deployment |
| `atlas-platform-chart/templates/ingress.yaml` | NGINX Ingress for `atlas-api.local` |

## Prerequisites

- Docker
- Minikube
- `kubectl`
- Helm
- Task CLI, if using `Taskfile.yaml`
- NGINX ingress controller for the ingress route

## Cluster Setup

From the `helm` directory, run:

```bash
task setup
```

This starts a four-node Minikube cluster.

The task also labels worker nodes for the chart scheduling rules:

| Node | Label | Used by |
|---|---|---|
| `minikube-m02` | `workload=app` | Atlas API Deployment |
| `minikube-m03` | `workload=db` | PostgreSQL StatefulSet |
| `minikube-m04` | `workload=dependent` | Reserved for dependent services |

If the cluster already exists, make sure these labels are present:

```bash
kubectl get nodes --show-labels
```

## Install The Chart

Install or upgrade the chart:

```bash
helm upgrade --install atlas-platform ./atlas-platform-chart
```

From the repository root, use:

```bash
helm upgrade --install atlas-platform ./helm/atlas-platform-chart
```

Check the deployed resources:

```bash
kubectl get all
kubectl get ingress
```

## Chart Resources

The chart renders these Kubernetes resources:

| Resource | Name | Purpose |
|---|---|---|
| `ConfigMap` | `atlas-cm` | Provides `ENV`, `PORT`, `VERSION`, `GOOSE_DRIVER`, and `GOOSE_MIGRATION_DIR` |
| `Secret` | `atlas-secrets` | Provides PostgreSQL and application database credentials |
| `Service` | `postgres` | Headless service for PostgreSQL |
| `StatefulSet` | `postgres` | Runs one PostgreSQL pod using `postgres:18` and a `512Mi` volume claim |
| `Job` | `atlas-migrate-job` | Waits for PostgreSQL and runs Goose migrations |
| `Service` | `atlas-api-svc` | Exposes the API inside the cluster on port `80` |
| `Deployment` | `atlas-api-deployment` | Runs two Atlas API replicas on port `3030` |
| `Ingress` | `atlas-api` | Routes `atlas-api.local/` to `atlas-api-svc` |

## Configuration

The default configurable image values are in `atlas-platform-chart/values.yaml`:

```yaml
images:
  api:
    repository: ankur2803/atlas-api
    tag: "v1"
    pullPolicy: IfNotPresent
  migrate:
    repository: ankur2803/atlas-migrate
    tag: "v1"
    pullPolicy: IfNotPresent
```

Override values at install time with:

```bash
helm upgrade --install atlas-platform ./atlas-platform-chart \
  --set images.api.tag=<tag> \
  --set images.migrate.tag=<tag>
```

Current template notes:

- The API image repository is hardcoded as `ghcr.io/ankurdubey28/atlas-api`.
- The migration image repository is hardcoded as `ghcr.io/ankurdubey28/atlas-migrate`.
- The API Deployment currently uses `{{ .Values.images.migrate.tag }}` for its image tag.
- The `pullPolicy` values are defined but are not currently referenced by the templates.

## Secrets

Use `templates/secrets.example.yaml` as the starting point for your own Kubernetes Secret manifest:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: atlas-secrets
type: Opaque
stringData:
  DB_URL: ""
  GOOSE_DBSTRING: ""
```

Fill in the values for your environment before installing the chart. Do not commit real credentials.

## Database Migrations

The migration Job:

1. Starts an init container using `postgres:alpine`.
2. Runs `pg_isready -h postgres -p 5432` until PostgreSQL is reachable.
3. Starts the `atlas-migrate` container.
4. Runs Goose with:

```bash
goose -dir "$GOOSE_MIGRATION_DIR" up
```

Check migration status with:

```bash
kubectl get jobs
kubectl logs job/atlas-migrate-job
```

If you need to rerun the migration Job, delete it first and upgrade the chart again:

```bash
kubectl delete job atlas-migrate-job
helm upgrade --install atlas-platform ./atlas-platform-chart
```

## Ingress

The chart creates an NGINX ingress for:

```text
http://atlas-api.local
```

For local Minikube usage, enable the ingress addon:

```bash
minikube addons enable ingress
```

Then map the host to your Minikube IP:

```bash
minikube ip
```

Add the returned IP to your hosts file:

```text
<minikube-ip> atlas-api.local
```

You can also port-forward the API service directly:

```bash
kubectl port-forward service/atlas-api-svc 3030:80
```

Then call:

```bash
curl http://localhost:3030/health
```

## Uninstall

Remove the Helm release:

```bash
helm uninstall atlas-platform
```

Delete the local Minikube cluster:

```bash
minikube delete
```
