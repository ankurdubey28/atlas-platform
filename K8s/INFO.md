# Atlas Platform Kubernetes Manifests

This directory contains the raw Kubernetes manifests for running Atlas Platform without Helm. The manifests deploy PostgreSQL, the API, database migrations, shared configuration, secrets, and ingress routing.

## Directory Layout

| Path | Purpose |
|---|---|
| `Taskfile.yml` | Helper tasks for local Minikube setup, manifest application, ingress setup, and cleanup |
| `manifest/cm.yaml` | ConfigMap for API and migration runtime configuration |
| `manifest/secrets.example.yaml` | Example Secret manifest to copy and fill with local values |
| `manifest/statefulset.yaml` | PostgreSQL headless Service and StatefulSet |
| `manifest/job.yaml` | Goose migration Job |
| `manifest/deployment.yaml` | Atlas API Service and Deployment |
| `manifest/ingress.yaml` | NGINX Ingress for `atlas-api.local` |

## Prerequisites

- Docker
- Minikube
- `kubectl`
- Task CLI, if using `Taskfile.yml`
- NGINX ingress controller for the ingress route

## Cluster Setup

From the `K8s` directory, run:

```bash
task spin-up-cluster
```

This starts a four-node Minikube cluster and labels worker nodes for the scheduling rules used by the manifests:

| Node | Label | Used by |
|---|---|---|
| `minikube-m02` | `workload=app` | Atlas API Deployment |
| `minikube-m03` | `workload=db` | PostgreSQL StatefulSet |
| `minikube-m04` | `workload=dependent` | Reserved for dependent services |

If the cluster already exists, confirm the labels:

```bash
kubectl get nodes --show-labels
```

## Secrets

Create your own Secret manifest from `manifest/secrets.example.yaml` and fill in the values for your environment before deployment.

The Secret must be named:

```text
atlas-secrets
```

The API, PostgreSQL, and migration Job all load environment variables from this Secret. Keep the Secret in the same namespace as the rest of the deployed resources. Do not commit real credentials.

## Deploy With Task

From the `K8s` directory, run:

```bash
task deploy
```

This runs the Taskfile deployment flow:

1. Start and label the Minikube cluster.
2. Apply the ConfigMap.
3. Apply the Secret.
4. Apply PostgreSQL.
5. Run database migrations.
6. Apply the API Service and Deployment.
7. Apply the Ingress.
8. Enable the Minikube ingress addon.
9. Port-forward the ingress controller to `localhost:8080`.

## Deploy Manually

Apply the manifests in dependency order:

```bash
kubectl apply -f manifest/cm.yaml
kubectl apply -f <your-secret-manifest>.yaml
kubectl apply -f manifest/statefulset.yaml
kubectl apply -f manifest/job.yaml
kubectl apply -f manifest/deployment.yaml
kubectl apply -f manifest/ingress.yaml
```

Enable ingress in Minikube:

```bash
minikube addons enable ingress
```

## Resources

The manifests create these Kubernetes resources:

| Resource | Name | Purpose |
|---|---|---|
| `ConfigMap` | `atlas-cm` | Provides `ENV`, `PORT`, `VERSION`, `GOOSE_DRIVER`, and `GOOSE_MIGRATION_DIR` |
| `Secret` | `atlas-secrets` | Provides application, migration, and PostgreSQL credentials |
| `Service` | `postgres` | Headless service for PostgreSQL |
| `StatefulSet` | `postgres` | Runs one PostgreSQL pod using `postgres:18` and a `512Mi` volume claim |
| `Job` | `atlas-migrate-job` | Waits for PostgreSQL and runs Goose migrations |
| `Service` | `atlas-api-svc` | Exposes the API inside the cluster on port `80` |
| `Deployment` | `atlas-api-deployment` | Runs two Atlas API replicas on port `3030` using `ankur2803/atlas-api:v1` |
| `Ingress` | `atlas-api` | Routes `atlas-api.local/` to `atlas-api-svc` |

## Runtime Configuration

`manifest/cm.yaml` provides these non-secret values:

| Key | Value | Used by |
|---|---|---|
| `ENV` | `PROD` | API runtime environment |
| `PORT` | `:3030` | API listen port |
| `VERSION` | `v1` | API version value |
| `GOOSE_DRIVER` | `postgres` | Migration database driver |
| `GOOSE_MIGRATION_DIR` | `/app/cmd/migrate/migrations` | Migration directory inside the container |

The container images are fixed in the manifests:

| Component | Image |
|---|---|
| API | `ankur2803/atlas-api:v1` |
| Migrations | `ankur2803/atlas-migrate:v1` |
| PostgreSQL | `postgres:18` |
| Migration readiness check | `postgres:alpine` |

## Database Migrations

The migration Job:

1. Starts an init container using `postgres:alpine`.
2. Runs `pg_isready -h postgres -p 5432` until PostgreSQL is reachable.
3. Starts the `atlas-migrate` container.
4. Runs Goose with:

```bash
goose -dir "$GOOSE_MIGRATION_DIR" up
```

Check migration status:

```bash
kubectl get jobs
kubectl logs job/atlas-migrate-job
```

Rerun migrations by deleting the Job and applying it again:

```bash
kubectl delete job atlas-migrate-job
kubectl apply -f manifest/job.yaml
```

## Ingress

The ingress routes this host:

```text
http://atlas-api.local
```

Map the host to your Minikube IP:

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

## Useful Checks

```bash
kubectl get all
kubectl get ingress
kubectl describe deployment atlas-api-deployment
kubectl describe statefulset postgres
kubectl logs deployment/atlas-api-deployment
```

## Cleanup

Delete the local Minikube cluster:

```bash
task destroy
```

Or run:

```bash
minikube delete
```
