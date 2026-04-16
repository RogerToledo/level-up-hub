#!/bin/bash

# Test runner script for Level Up Hub backend
# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║  Level Up Hub - Test Suite Runner     ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""

# Function to run tests for a package
run_test() {
    local package=$1
    local name=$2
    
    echo -e "${YELLOW}→ Running tests for ${name}...${NC}"
    
    if go test "./${package}" -v -coverprofile="${package////_}_coverage.out" 2>&1 | grep -E "PASS|FAIL|coverage:"; then
        echo -e "${GREEN}✓ ${name} tests passed${NC}"
        echo ""
        return 0
    else
        echo -e "${RED}✗ ${name} tests failed${NC}"
        echo ""
        return 1
    fi
}

# Track failures
FAILED_TESTS=""
EXIT_CODE=0

# Run all test suites
echo -e "${BLUE}═══ Core Components ═══${NC}"
run_test "internal/pagination" "Pagination" || { FAILED_TESTS="${FAILED_TESTS}\n- Pagination"; EXIT_CODE=1; }
run_test "internal/logger" "Logger" || { FAILED_TESTS="${FAILED_TESTS}\n- Logger"; EXIT_CODE=1; }
run_test "internal/pkg/identity" "Identity" || { FAILED_TESTS="${FAILED_TESTS}\n- Identity"; EXIT_CODE=1; }

echo -e "${BLUE}═══ Business Logic ═══${NC}"
run_test "internal/activity" "Activity" || { FAILED_TESTS="${FAILED_TESTS}\n- Activity"; EXIT_CODE=1; }
run_test "internal/ladder" "Ladder" || { FAILED_TESTS="${FAILED_TESTS}\n- Ladder"; EXIT_CODE=1; }

echo -e "${BLUE}═══ Infrastructure ═══${NC}"
run_test "internal/email" "Email" || { FAILED_TESTS="${FAILED_TESTS}\n- Email"; EXIT_CODE=1; }
run_test "internal/rest" "REST" || { FAILED_TESTS="${FAILED_TESTS}\n- REST"; EXIT_CODE=1; }
run_test "internal/database" "Database" || { FAILED_TESTS="${FAILED_TESTS}\n- Database"; EXIT_CODE=1; }
run_test "apperr" "Error Handling" || { FAILED_TESTS="${FAILED_TESTS}\n- Error Handling"; EXIT_CODE=1; }

echo -e "${BLUE}═══════════════════════${NC}"
echo ""

# Generate combined coverage report
echo -e "${YELLOW}→ Generating combined coverage report...${NC}"
echo "mode: set" > coverage.out
for file in *_coverage.out; do
    if [ -f "$file" ]; then
        grep -h -v "mode: set" "$file" >> coverage.out
    fi
done

COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
echo -e "${GREEN}Total Coverage: ${COVERAGE}${NC}"
echo ""

# Summary
if [ $EXIT_CODE -eq 0 ]; then
    echo -e "${GREEN}╔════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║    ✓ All tests passed successfully!   ║${NC}"
    echo -e "${GREEN}╚════════════════════════════════════════╝${NC}"
else
    echo -e "${RED}╔════════════════════════════════════════╗${NC}"
    echo -e "${RED}║         ✗ Some tests failed:          ║${NC}"
    echo -e "${RED}║${FAILED_TESTS}${NC}"
    echo -e "${RED}╚════════════════════════════════════════╝${NC}"
fi

# Open coverage report in browser if requested
if [ "$1" == "--coverage" ]; then
    echo -e "${YELLOW}→ Opening coverage report in browser...${NC}"
    go tool cover -html=coverage.out -o coverage.html
    if command -v open &> /dev/null; then
        open coverage.html
    elif command -v xdg-open &> /dev/null; then
        xdg-open coverage.html
    fi
fi

# Run benchmarks if requested
if [ "$1" == "--bench" ]; then
    echo -e "${BLUE}═══ Running Benchmarks ═══${NC}"
    go test -bench=. -benchmem ./internal/activity ./internal/email ./internal/rest ./apperr
fi

exit $EXIT_CODE
