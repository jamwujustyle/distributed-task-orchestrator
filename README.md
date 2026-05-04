# Distributed Task Orchestrator (Go + AWS + LocalStack)

A high-performance, distributed system designed to manage, store, and execute scripts across a decoupled architecture. Built with Go, leveraging AWS S3 for artifact storage and DynamoDB for state management, with LocalStack providing the development environment.

## 🏗 System Architecture

The system is split into two primary binaries to ensure scalability and separation of concerns:

- **Controller** (`cmd/controller`): The "Brain."
  - Exposes a gRPC/API for task submission.
  - Uploads script artifacts to S3.
  - Initializes task state in DynamoDB as `PENDING`.

- **Worker** (`cmd/worker`): The "Muscle."
  - Polls/Receives tasks.
  - Downloads scripts from S3.
  - Executes the logic and updates DynamoDB status (`RUNNING` ➔ `COMPLETED`/`FAILED`).

## 🛠 Tech Stack

- **Language:** Go (Golang)
- **Infrastructure:** LocalStack (Simulating AWS S3 & DynamoDB)
- **Logging:** Custom `slog` wrapper with colored terminal output and source-line tracking.
- **Storage Layer:**
  - **S3:** Content-addressable storage for task scripts.
  - **DynamoDB:** Low-latency NoSQL for task metadata and lifecycle tracking.

## 📂 Project Structure

```plaintext
.
├── cmd
│   ├── controller    # Main entry point for the API/Controller
│   ├── worker        # Main entry point for the Execution Worker
│   └── cli           # Admin CLI tool for manual task management
├── internal
│   ├── store         # Storage Layer (S3 & DynamoDB logic)
│   ├── logger        # Custom slog initialization and handlers
│   └── engine        # Business logic for task execution
├── pkg
│   └── protocol      # Protobuf definitions and generated gRPC code
└── Makefile          # Automation for LocalStack and Builds
```

## 💡 Key Design Decisions

### 1. Lazy Connections & Startup Pings

In Go, AWS clients are initialized as structural shells. To ensure the environment is ready, we implement a `Ping()` method using `HeadBucket` (S3) and `DescribeTable` (DynamoDB). This forces a "fail-fast" behavior if LocalStack is offline.

### 2. Dependency Injection

We pass `aws.Config` into store constructors (`NewS3Store`, `NewDynamoStore`). This allows us to point the entire system to LocalStack for development or the real AWS Cloud for production just by changing environment variables.

### 3. Atomic Updates

We use Expression Attribute Names (`#status`) in DynamoDB updates to avoid conflicts with reserved words and ensure state transitions are consistent across multiple workers.

## 🚀 Development Workflow

### Infrastructure Setup

Ensure LocalStack is running via Docker, then initialize the storage:

```bash
# Create S3 Bucket
awslocal s3 mb s3://tasks-bucket

# Create DynamoDB Table
awslocal dynamodb create-table \
  --table-name TasksTable \
  --attribute-definitions AttributeName=ID,AttributeType=S \
  --key-schema AttributeName=ID,KeyType=HASH \
  --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5
```

### Running the App

```bash
# Set local environment variables
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test

# Build and run
make run
```

## 📋 Current Data Models

### Task Object

| Field         | Type     | Description                                  |
| ------------- | -------- | -------------------------------------------- |
| `ID`          | `string` | Unique UUID for the task.                    |
| `ScriptS3Key` | `string` | Path to the script file in S3.               |
| `Status`      | `enum`   | `PENDING`, `RUNNING`, `COMPLETED`, `FAILED`. |
| `CreatedAt`   | `int64`  | Unix timestamp of creation.                  |
| `UpdatedAt`   | `int64`  | Unix timestamp of last state change.         |

### GOAL

GAIN MASTERY OF THE STACK
