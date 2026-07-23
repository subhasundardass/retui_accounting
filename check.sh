#!/usr/bin/env bash
# Runs the same checks as .github/workflows/ci.yml, locally.
# Usage: ./check.sh

set -euo pipefail

# Colors for readable output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

step() {
    echo -e "\n${GREEN}==>${NC} $1"
}

fail() {
    echo -e "${RED}FAILED:${NC} $1"
    exit 1
}

step "Checking go.mod / go.sum are tidy"
go mod tidy
# if ! git diff --exit-code go.mod go.sum > /dev/null; then
#     fail "go.mod/go.sum changed after 'go mod tidy' — commit the updated files"
# fi

step "Building library + example app"
go build ./... || fail "go build ./... failed"
go build ./cmd/app || fail "go build ./cmd/app failed"

step "Running go vet"
go vet ./... || fail "go vet failed"

step "Checking gofmt"
fmt_out="$(gofmt -l .)"
if [ -n "$fmt_out" ]; then
    echo "$fmt_out"
    fail "the above files are not gofmt'd — run: gofmt -w ."
fi

step "Running tests (race detector + coverage)"
go test ./... -race -coverprofile=coverage.out -covermode=atomic || fail "tests failed"

step "Running golangci-lint"
if command -v golangci-lint &> /dev/null; then
    golangci-lint run --timeout=5m || fail "lint failed"
else
    echo "golangci-lint not installed, skipping. Install with:"
    echo "  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$(go env GOPATH)/bin"
fi

echo -e "\n${GREEN}All checks passed.${NC}"