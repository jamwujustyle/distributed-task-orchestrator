# Distributed Task Orchestrator

A production-style distributed system for submitting, queuing, and executing sandboxed scripts across a multi-node Kubernetes cluster. Built from scratch as a learning project to gain hands-on mastery of backend infrastructure, cloud-native tooling, and distributed systems patterns.

## Architecture Overview

```
CLI ──gRPC──▸ Controller ──Kafka──▸ Workers ──AWS SDK(localstack)──▸ Lambda (Sandboxed Execution)
                │                      │                        │
                ▼                      ▼                        ▼
             DynamoDB              S3 (Scripts)           Restricted Shell
           (Task State)          (Artifact Store)      (Whitelisted Binaries)
```

**Controller** — gRPC API that accepts task submissions, uploads scripts to S3, publishes events to Kafka, and tracks task lifecycle in DynamoDB.

**Workers** — Kafka consumers that pick up tasks, retrieve scripts from S3, invoke Lambda for sandboxed execution, and report results back to DynamoDB.

**Lambda Sandbox** — Executes user scripts under a 4-layer security model: static analysis, environment cleansing, restricted `PATH` via symlinked binaries, and OS-level privilege dropping (`nobody` user).

## Tech Stack

| Layer | Technologies |
| :--- | :--- |
| **Language** | Go |
| **Communication** | gRPC, Protocol Buffers, Kafka (event-driven messaging) |
| **Cloud Services** | AWS Lambda, S3, DynamoDB (via LocalStack) |
| **Orchestration** | Kubernetes (K3s), multi-node Vagrant VMs |
| **Provisioning** | Ansible (playbooks + Helm automation) |
| **Observability** | Prometheus, Grafana, Loki + Promtail, OpenTelemetry, Jaeger |
| **Security** | Lambda sandboxing (command whitelisting, env stripping, privilege drop) |
| **Tooling** | Docker, Just, Air (hot-reload) |

## Infrastructure

- **Vagrant** spins up multi-node VMs (control plane + worker nodes)
- **Ansible** provisions each node end-to-end: installs Docker, bootstraps K3s, joins worker nodes to the cluster, and deploys services via Helm charts (Kafka, Prometheus, Loki stack)
- **Kubernetes manifests** (`k8s/`) define all application deployments, services, config, and secrets — including a LocalStack pod with an init container that extracts and registers the Lambda binary on startup

## Local Development

- **Docker Compose** runs the full stack locally (LocalStack, Kafka, controller, worker) for rapid iteration
- **Air** provides hot-reload on code changes — saves are reflected instantly without manual rebuilds
- **CLI** (`cmd/cli`) connects to the controller via gRPC for uploading scripts, submitting tasks, and retrieving results in real time

## Purpose

This project exists purely for learning. The goal is end-to-end fluency across the modern backend stack — not just writing Go, but wiring together messaging, infrastructure-as-code, container orchestration, observability, and security into a single cohesive system.
