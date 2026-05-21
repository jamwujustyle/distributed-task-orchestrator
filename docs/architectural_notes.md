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

### Next Steps & Backlog

#### 1. Observability: Metrics & Monitoring
* **Problem:** We have centralized logs, but zero visibility into metrics.
* **Solution:**
  * Instrument the Go controller and worker codebases with Prometheus metrics (e.g., tracking task count, execution latency, processing rates).
  * Deploy a Prometheus server inside the K3s cluster to scrape metrics endpoints.
  * Construct Grafana dashboards combining Loki logs and Prometheus metrics.

#### 2. Distributed Tracing (OpenTelemetry)
* **Problem:** Tracking a single task's flow across CLI -> Controller -> Kafka -> Worker -> Lambda -> DynamoDB is difficult when diagnosing failure bottlenecks.
* **Solution:** Implement context propagation using OpenTelemetry trace spans across the gRPC and Kafka network boundaries.

#### 3. Lambda script execution Sandboxing
* **Problem:** Workers invoke AWS Lambda functions executing raw bash scripts without validation, which poses security risks.
* **Solution:** Apply command whitelisting, configure memory/CPU constraints, and namespace executions.

#### 4. Integration Tests
* **Problem:** End-to-end task flows rely on manual testing via the CLI.
* **Solution:** Create integration test suites to simulate and assert task lifecycle flows (Kafka publishing, LocalStack execution, S3 storage, and DynamoDB updates).
