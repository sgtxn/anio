lint:
	golangci-lint run -c .golangci.yml

lint_fix:
	golangci-lint run -c .golangci.yml --fix
