.PHONY: run
run:
	go run ./cmd/web

.PHONY: dev
dev:
	ls **/*.* | entr -cr go run ./cmd/web

.PHONY: tidy
tidy:
	@echo 'Formatting .go files...'
	go fmt ./...
	@echo 'Tidying module dependencies...'
	go mod tidy
	@echo 'Verifying and vendoring module dependencies...'
	go mod verify
	go mod vendor

.PHONY: audit
audit:
	@echo 'Checking module dependencies'
	go mod tidy -diff
	go mod verify
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

.PHONY: build
build:
	go build ./cmd/projectspy

.PHONY: install
install:
	mkdir -p $$HOME/bin
	mv projectspy $$HOME/bin/pspy
