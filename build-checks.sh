#!/bin/bash
COVERAGE_PERCENT=0
echo "Running build checks..."
# Check that tests pass
gotestsum --junitfile ./test-reports/junit.xml -- --coverprofile=coverage.out ./...
RC=$?
if [[ $RC -ne 0 ]]; then
    echo "Test have failed. Please fix your tests and try committing again."
    exit 1
fi
echo "Tests have passed!"

# Check that test coverage is sufficient
echo "Checking code coverage..."
COVERAGE=$(go tool cover -func=coverage.out | \
    tail -n 1 | \
    awk '{ print $3 }' | \
    sed -e 's/^\([0-9]*\).*$/\1/g')
if [ "$CI" == "true" ]; then
    if [ "$COVERAGE" -lt "$COVERAGE_PERCENT" ]; then COVERAGE_RESULT=FAILED; else COVERAGE_RESULT=PASSED; fi
    curl -s --proxy 'http://localhost:29418' --request PUT "http://api.bitbucket.org/2.0/repositories/${BITBUCKET_REPO_OWNER}/${BITBUCKET_REPO_SLUG}/commit/${BITBUCKET_COMMIT}/reports/go-test-cover" \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "title": "Test coverage",
        "details": "Test coverage from go test",
        "reporter": "go test",
        "report_type": "COVERAGE",
        "result": "'"${COVERAGE_RESULT}"'",
        "data": [
            {
                "title": "Statement coverage",
                "type": "PERCENTAGE",
                "value": '"${COVERAGE}"'
            },
            {
                "title": "Required coverage",
                "type": "PERCENTAGE",
                "value": '"${COVERAGE_PERCENT}"'
            }
        ]
    }'
fi
if [ "$COVERAGE" -lt "$COVERAGE_PERCENT" ]; then
    echo "Code coverage is currently at $COVERAGE%. This needs to be at least $COVERAGE_PERCENT% before being able to commit."
    exit 1
fi
echo "Coverage is $COVERAGE%. Nice!"

# Run the linter
echo "Running linter..."
echo "If you don't have golang-lint ci installed, just run the following command:"
echo "curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.47.2"

golangci-lint run --exclude-use-default=false
RC=$?
if [[ RC -ne 0 ]]; then
    echo "Linter has found issues. Please resolve these before committing your changes."
    exit 1
fi
# Linter will return exit code 1 if it finds any issues
echo "Linter found no issues."
echo -e '\(^-^)/ Build checks are all good!' 
