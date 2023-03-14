default: build

get:
	go get ./...

test:
	gotestsum --junitfile ./test-reports/junit.xml -- --coverprofile=coverage.out ./src/...
	go tool cover -func coverage.out

buildchecks:
	bash ./build-checks.sh
