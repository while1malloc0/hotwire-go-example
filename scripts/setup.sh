#!/usr/bin/env bash

cd "$(dirname "$0")/.." || exit 1

./scripts/bootstrap.sh

echo "Successfully setup project. Run the server with 'go run main.go'"