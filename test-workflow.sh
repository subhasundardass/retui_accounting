#!/bin/bash

echo "Testing CI workflow locally for retui packages..."

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Test 1: Build
echo "📦 Testing build..."
if go build ./retui/...; then
    echo -e "${GREEN}✅ Build passed${NC}"
else
    echo -e "${RED}❌ Build failed${NC}"
    exit 1
fi

# Test 2: Vet
echo "🔍 Testing vet..."
if go vet ./retui/...; then
    echo -e "${GREEN}✅ Vet passed${NC}"
else
    echo -e "${RED}❌ Vet failed${NC}"
    exit 1
fi

# Test 3: Gofmt
echo "📝 Testing gofmt..."
if [ -z "$(gofmt -l ./retui)" ]; then
    echo -e "${GREEN}✅ Gofmt passed${NC}"
else
    echo -e "${RED}❌ Gofmt failed${NC}"
    exit 1
fi

# Test 4: Tests with coverage
echo "🧪 Testing tests..."
if go test ./retui/... -v -race -coverprofile=coverage.out; then
    echo -e "${GREEN}✅ Tests passed${NC}"
    go tool cover -func=coverage.out | grep total
else
    echo -e "${RED}❌ Tests failed${NC}"
    exit 1
fi

# Test 5: Lint with golangci-lint
echo "🔎 Testing lint (golangci-lint)..."

# Install golangci-lint if not present
if ! command_exists golangci-lint; then
    echo -e "${YELLOW}⚠️  golangci-lint not found. Installing...${NC}"
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
fi

if golangci-lint run ./retui/... --timeout=5m; then
    echo -e "${GREEN}✅ golangci-lint passed${NC}"
else
    echo -e "${RED}❌ golangci-lint failed${NC}"
    exit 1
fi



# Test 9: govulncheck (optional)
echo "🔒 Testing govulncheck..."

if ! command_exists govulncheck; then
    echo -e "${YELLOW}⚠️  govulncheck not found. Installing...${NC}"
    go install golang.org/x/vuln/cmd/govulncheck@latest
fi

if govulncheck ./retui/... 2>/dev/null; then
    echo -e "${GREEN}✅ govulncheck passed${NC}"
else
    echo -e "${YELLOW}⚠️  govulncheck found issues (not failing)${NC}"
fi

# Test 10: Dependency check
echo "📦 Checking dependencies..."
if go mod tidy && git diff --exit-code go.mod go.sum 2>/dev/null; then
    echo -e "${GREEN}✅ Dependencies are tidy${NC}"
else
    echo -e "${RED}❌ Dependencies are not tidy. Run 'go mod tidy'${NC}"
    exit 1
fi

echo -e "${GREEN}✅ All tests passed!${NC}"
echo -e "${GREEN}📊 Coverage report saved to coverage.out${NC}"