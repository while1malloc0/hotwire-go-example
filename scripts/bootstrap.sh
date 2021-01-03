#!/usr/bin/env bash

cd "$(dirname "$0")/.." || exit 1

set -euo pipefail

echo "Checking for node"
command -v npm &> /dev/null || { echo "npm not found. Install it and try again."; exit 1; }

echo "Checking for go"
command -v go &> /dev/null || { echo "go not found. Install it and try again."; exit 1; }

echo "Installing Go deps"
go mod download 1>/dev/null 2>err.log || { echo "Failed to install Go deps. Saved error to err.log"; cat err.log; exit 1; }

echo "Installing node deps"
npm install 1>/dev/null 2>err.log || { echo "Failed to install node deps. Saved error to err.log"; cat err.log; exit 1; }

echo "Creating webpack bundle"
npx webpack -c webpack.config.js --mode development 1>/dev/null 2>err.log || { echo "Failed to created webpack bundle. Saved error to err.log"; cat err.log; exit 1; }

rm err.log &> /dev/null

echo "Successfully bootstrapped project dependencies"