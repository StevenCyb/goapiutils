install:
	echo "Install:"
	go mod download
	go install mvdan.cc/gofumpt@latest
	go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest

format:
	echo "Format:"
	gofmt -w -e .
	gofumpt -l -w .
	find . -type f -name "*.go" -execdir fieldalignment -fix {} \;
