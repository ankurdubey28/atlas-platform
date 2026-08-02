# Argo CD Local Setup

This guide installs Argo CD on a local Kubernetes cluster such as Minikube.

## Prerequisites

* Docker
* `kubectl`
* Minikube or another local Kubernetes cluster
* A running Kubernetes cluster

Verify the cluster:

```bash
kubectl cluster-info
kubectl get nodes
```

## 1. Create the Argo CD namespace

```bash
kubectl create namespace argocd
```

## 2. Install Argo CD

From a cloned Argo CD repository:

```bash
kubectl apply \
  -n argocd \
  --server-side \
  --force-conflicts \
  -f manifests/install.yaml
```

## 3. Check the installation

```bash
kubectl get pods -n argocd
```

Wait until all Argo CD pods are running.

## 4. Get the admin password

```bash
kubectl -n argocd get secret argocd-initial-admin-secret \
  -o jsonpath="{.data.password}" | base64 -d; echo
```

The default username is:

```text
admin
```

## 5. Access the Argo CD UI

Forward the Argo CD server port:

```bash
kubectl port-forward service/argocd-server \
  -n argocd 8080:443
```

Open:

```text
https://localhost:8080
```

Your browser may display a certificate warning because Argo CD uses a self-signed certificate locally.

## Optional: Set the default namespace

```bash
kubectl config set-context --current --namespace=argocd
```
