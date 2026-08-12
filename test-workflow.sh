#!/bin/bash

echo "Testing CI workflow locally for retui packages..."

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Log files
LOG_FILE="ci-test.log"

# Start fresh log
: > "$LOG_FILE"

echo "==================================================" >> "$LOG_FILE"
echo "RetUI CI Test Log" >> "$LOG_FILE"
echo "Started: $(date)" >> "$LOG_FILE"
echo "==================================================" >> "$LOG_FILE"

# Pin versions to match .github/workflows/ci.yml
GOLANGCI_LINT_VERSION="v2.12.2"
GOSEC_VERSION="latest"

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Log a failed test
log_failure() {
    local test_name="$1"

    echo "" >> "$LOG_FILE"
    echo "==================================================" >> "$LOG_FILE"
    echo "FAILED: $test_name" >> "$LOG_FILE"
    echo "Time: $(date)" >> "$LOG_FILE"
    echo "==================================================" >> "$LOG_FILE"
}

# Track overall failure
FAILED=0

# ==================================================
# Test 1: Build
# ==================================================

echo "📦 Testing build (native)..."

if go build ./retui/... 2>&1 | tee -a "$LOG_FILE"; then
    echo -e "${GREEN}✅ Build passed (native)${NC}"
else
    echo -e "${RED}❌ Build failed (native)${NC}"
    log_failure "Build (native)"
    echo -e "${RED}📄 See $LOG_FILE for details${NC}"
    exit 1
fi

# ==================================================
# Test 1b: Windows Build
# ==================================================

echo "📦 Testing build (windows)..."

if GOOS=windows GOARCH=amd64 go build ./retui/... 2>&1 | tee -a "$LOG_FILE"; then
    echo -e "${GREEN}✅ Build passed (windows)${NC}"
else
    echo -e "${RED}❌ Build failed (windows)${NC}"
    log_failure "Build (windows)"
    echo -e "${RED}📄 See $LOG_FILE for details${NC}"
    exit 1
fi

# ==================================================
# Test 2: Vet
# ==================================================

echo "🔍 Testing vet..."

if go vet ./retui/... 2>&1 | tee -a "$LOG_FILE"; then
    echo -e "${GREEN}✅ Vet passed${NC}"
else
    echo -e "${RED}❌ Vet failed${NC}"
    log_failure "go vet"
    echo -e "${RED}📄 See $LOG_FILE for details${NC}"
    exit 1
fi

# ==================================================
# Test 3: Gofmt
# ==================================================

echo "📝 Testing gofmt..."

fmt_out="$(gofmt -l ./retui)"

if [ -z "$fmt_out" ]; then
    echo -e "${GREEN}✅ Gofmt passed${NC}"
else
    echo -e "${RED}❌ Gofmt failed${NC}"

    log_failure "gofmt"

    echo "The following files are not gofmt'd:" | tee -a "$LOG_FILE"
    echo "$fmt_out" | tee -a "$LOG_FILE"

    echo -e "${RED}📄 See $LOG_FILE for details${NC}"
    exit 1
fi

# ==================================================
# Test 4: Tests + Coverage
# ==================================================
echo "🧪 Testing tests..."

TEST_OUTPUT=$(mktemp)

if go test ./retui/... -v -race \
    -covermode=atomic \
    -coverprofile=coverage.out \
    -coverpkg=./retui/... 2>&1 | tee "$TEST_OUTPUT"; then

    echo -e "${GREEN}✅ Tests passed${NC}"
    go tool cover -func=coverage.out | grep total
    rm -f "$TEST_OUTPUT"

else
    echo -e "${RED}❌ Tests failed${NC}"

    # Keep only failure-related lines
    grep -E '^(--- FAIL:|FAIL|    .*:.*|        .*)' \
        "$TEST_OUTPUT" > fail-test.log

    echo -e "${RED}❌ Failed test details saved to fail-test.log${NC}"
    cat fail-test.log

    rm -f "$TEST_OUTPUT"
    exit 1
fi

# ==================================================
# Test 5: golangci-lint
# ==================================================

# echo "🔎 Testing lint (golangci-lint ${GOLANGCI_LINT_VERSION})..."

# installed_version=""

# if command_exists golangci-lint; then
#     installed_version="$(golangci-lint version 2>/dev/null | grep -oE 'version [^ ]+' | awk '{print $2}')"
# fi

# if [ "$installed_version" != "${GOLANGCI_LINT_VERSION#v}" ]; then
#     echo -e "${YELLOW}⚠️  golangci-lint ${GOLANGCI_LINT_VERSION} not found. Installing...${NC}"
#     go install "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@${GOLANGCI_LINT_VERSION}"
# fi

# if golangci-lint run ./retui/... --timeout=5m 2>&1 | tee -a "$LOG_FILE"; then
#     echo -e "${GREEN}✅ golangci-lint passed${NC}"
# else
#     echo -e "${RED}❌ golangci-lint failed${NC}"
#     log_failure "golangci-lint"
#     echo -e "${RED}📄 See $LOG_FILE for details${NC}"
#     exit 1
# fi

# ==================================================
# Test 6: gosec
# ==================================================

echo "🛡️  Testing gosec..."

if ! command_exists gosec; then
    echo -e "${YELLOW}⚠️  gosec not found. Installing...${NC}"
    go install github.com/securego/gosec/v2/cmd/gosec@latest
fi

if gosec ./retui/... 2>&1 | tee -a "$LOG_FILE"; then
    echo -e "${GREEN}✅ gosec passed${NC}"
else
    echo -e "${RED}❌ gosec failed${NC}"
    log_failure "gosec"
    echo -e "${RED}📄 See $LOG_FILE for details${NC}"
    exit 1
fi

# ==================================================
# Test 7: govulncheck
# ==================================================

echo "🔒 Testing govulncheck..."

if ! command_exists govulncheck; then
    echo -e "${YELLOW}⚠️  govulncheck not found. Installing...${NC}"
    go install golang.org/x/vuln/cmd/govulncheck@latest
fi

if govulncheck ./retui/... 2>&1 | tee -a "$LOG_FILE"; then
    echo -e "${GREEN}✅ govulncheck passed${NC}"
else
    echo -e "${YELLOW}⚠️  govulncheck found issues (not failing)${NC}"

    log_failure "govulncheck (non-blocking)"
fi

# ==================================================
# Test 8: Dependency Check
# ==================================================

echo "📦 Checking dependencies..."

go mod tidy 2>&1 | tee -a "$LOG_FILE"

if git diff --exit-code go.mod go.sum >> "$LOG_FILE" 2>&1; then
    echo -e "${GREEN}✅ Dependencies are tidy${NC}"
else
    echo -e "${RED}❌ Dependencies are not tidy${NC}"

    log_failure "Dependency check"

    echo "Run 'go mod tidy' and commit the changes." | tee -a "$LOG_FILE"

    echo -e "${RED}📄 See $LOG_FILE for details${NC}"
    exit 1
fi

# ==================================================
# Success
# ==================================================

echo "" >> "$LOG_FILE"
echo "==================================================" >> "$LOG_FILE"
echo "ALL TESTS PASSED" >> "$LOG_FILE"
echo "Finished: $(date)" >> "$LOG_FILE"
echo "==================================================" >> "$LOG_FILE"

echo -e "${GREEN}✅ All tests passed!${NC}"
echo -e "${GREEN}📊 Coverage report saved to coverage.out${NC}"
echo -e "${GREEN}📄 Test log saved to $LOG_FILE${NC}"