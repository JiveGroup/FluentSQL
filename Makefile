.PHONY: critic security vulncheck lint test all

mod:
	go list -m --versions

test:
	go test -v -timeout 30s -coverprofile=coverage.txt -cover ./...
	go tool cover -func=coverage.txt

critic:
	gocritic check -enableAll -disable=unnamedResult,unlabelStmt,hugeParam,singleCaseSwitch,builtinShadow,typeAssertChain ./...

security:
	gosec -exclude-dir=mysql,psql -exclude=G103,G115,G401,G501 ./...

vulncheck:
	govulncheck ./...

lint:
	golangci-lint run ./...

check: critic security vulncheck lint

all: critic security vulncheck lint test