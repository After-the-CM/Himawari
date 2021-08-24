.PHONY: build
## build: build the "Himawari"
build: clean
	@echo "Building..."
	@go build

.PHONY: run
## run: runs Himawari and mysql container
run:build
	./Himawari

.PHONY: clean
## clean: cleans the binary
clean:
	@echo "Cleaning"
	@go clean

.PHONY: setup
## setup: setup go modules
setup:
	@go mod init Himawari\
		&& go mod tidy
	
.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
