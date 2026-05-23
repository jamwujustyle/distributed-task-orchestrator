A Critical Note on your Lambda Code
Look at this line in the code Claude gave you:
cmd := exec.CommandContext(ctx, tmp.Name())

This is a security nightmare waiting to happen. You are taking a Base64 string from the internet, writing it to a file, and executing it as a shell script.

If this were production: Someone could send rm -rf / or a crypto miner.

Since it's LocalStack: It’s "safe" for testing, but as you move toward a "High-Value Engineer" profile, you'll want to think about sandboxing.

Architecture Shift
In a real orchestrator, the Lambda shouldn't just run raw bash. It should:

Validate the script against a whitelist of commands.

Limit the execution time (which Lambda does via timeout).

Namespace the execution so it can't see the Lambda's environment variables (like your AWS keys).

<!--  -->
<!--  -->
<!-- PLAN -->

Finish this app
↓
Deploy to real AWS (EC2 + real DynamoDB/S3)
↓
Add Ansible playbooks to automate deployment
↓
Migrate to EKS (managed K8s on AWS)
↓
Add Kafka for task queue instead of gRPC polling

Next options:

Kafka — replace gRPC polling with event-driven task queue. Worker subscribes to a topic, controller publishes when task is submitted.
Ansible — write playbooks to provision and deploy your stack to a real server.
Observability — add metrics/tracing with OpenTelemetry.
Tests — write integration tests for the service layer.

Given your roadmap you wrote earlier:
Finish app ✓
Deploy to AWS
Add Ansible ✓
Migrate to EKS ✓ (K3s clustering on Vagrant VMs)
Add Kafka ✓

### Current State of the Architecture (As of May 2026)
* **Event-Driven Messaging:** Successfully migrated from gRPC polling to Kafka. The controller acts as a producer to the `tasks` topic, and workers consume messages from it.
* **Kubernetes (K3s):** Shifted infrastructure deployment from Docker Compose to a K3s cluster. Deployments are orchestrated across a multi-node Vagrant VM environment.
* **Ansible Infrastructure Provisioning:** Ansible playbooks automate the setup of the VM dependencies, Docker, and K3s node clustering.
* **Centralized Logging:** Deployed `loki-stack` (Loki, Promtail, and Grafana) to the Kubernetes cluster. Promtail automatically scrapes stdout/stderr logs from all running pods.
* **Distributed Observability (Metrics & Monitoring):** Instrumented Go controller and worker codebases to serve Prometheus metrics (tracking tasks submitted, tasks processed, and task execution durations) on port `9091`. Deployed Prometheus server using Ansible Helm automation, enabling real-time metrics dashboards in Grafana.
* **Distributed Tracing (OpenTelemetry):** Fully integrated OpenTelemetry tracing across the entire system. Traces capture the full end-to-end task lifecycle (CLI ➔ gRPC ➔ Controller ➔ Kafka ➔ Worker ➔ Status Updates), maintaining trace-context propagation across network and message boundaries. Traces are gathered and visualized using a Jaeger all-in-one backend inside the cluster.
* **Lambda Sandboxing:** Implemented a 4-layer defense-in-depth security model for script execution: (1) static analysis blocking absolute/relative path bypasses, (2) environment cleansing stripping AWS credentials from the child process, (3) restricted `PATH` via symlinked whitelisted binaries (`echo`, `cat`, `grep`, `sleep`, `date`), and (4) OS-level privilege dropping to the `nobody` user (UID 65534).

### Next Steps & Backlog

#### 1. Integration Tests
* **Problem:** End-to-end task flows rely on manual testing via the CLI.
* **Solution:** Create integration test suites to simulate and assert task lifecycle flows (Kafka publishing, LocalStack execution, S3 storage, and DynamoDB updates).
