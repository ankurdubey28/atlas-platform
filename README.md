# atlas-platform

**atlas-platform** is a production-oriented infrastructure project built to simulate how real-world services are developed, deployed, and operated.

The goal of this project is to understand the complete lifecycle of a backend service — starting from local development to production deployment — while focusing on reliability, automation, and observability.

---

## What is a Platform

A platform is a combination of infrastructure, tools, and abstractions used to build, deploy, and operate applications efficiently.

---

## What this Project Covers

- Building a REST API (Go)
- Containerization using Docker
- Local development setup
- CI/CD pipeline setup
- Deployment on Kubernetes
- GitOps-based delivery using ArgoCD
- Observability (metrics, logs, alerts)
- Debugging production-like scenarios

---

## Project Structure

```
atlas-platform/
├── api/                # Go REST API
├── docker/             # Docker related files
├── k8s/                # Kubernetes manifests
├── helm/               # Helm charts (to be added)
├── argocd/             # GitOps configs (to be added)
├── ci/                 # CI pipeline configs
└── observability/      # Monitoring setup (to be added)
```

---

## Progress 

- [x] REST API server
- [x] Docker containerization
- [ ] Local development setup
- [ ] CI pipeline
- [ ] Bare metal deployment
- [ ] Kubernetes cluster setup
- [ ] Kubernetes deployment
- [ ] Helm charts
- [ ] ArgoCD setup
- [ ] Observability stack
- [ ] Dashboards & alerts

---

## Tech Stack

- Language: Go
- Containerization: Docker
- Orchestration: Kubernetes
- CI/CD: (to be added)
- GitOps: ArgoCD (planned)
- Observability: Prometheus, Grafana (planned)

---

## Note

This project is being built incrementally as part of an SRE learning journey.  
The README will be updated as new components are added.