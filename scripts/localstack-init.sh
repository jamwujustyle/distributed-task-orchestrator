#!/usr/bin/env bash

set -euo pipefail

REGION=us-east-1

echo "Creating a new s3 bucket: tasks-bucket"
awslocal s3api create-bucket --bucket tasks-bucket --region $REGION


echo "Creating a new dynamodb: TaskTable"
awslocal dynamodb create-table \
    --table-name TaskTable \
    --attribute-definitions AttributeName=ID,AttributeType=S \
    --key-schema AttributeName=ID,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --region $REGION

echo "localstack initialized"
