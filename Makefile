.PHONY: test test-coverage test-package test-user open-coverage

# Menjalankan semua test
test:
	go test -v ./...

# Menjalankan test dengan coverage dan generate HTML report
test-coverage:
	@mkdir -p tmp
	go test -coverprofile=tmp/coverage.out ./...
	go tool cover -html=tmp/coverage.out -o tmp/coverage.html
	@echo "Coverage report generated at tmp/coverage.html"
	@open tmp/coverage.html

# Menjalankan test untuk package tertentu dengan coverage
test-package:
	@if [ -z "$(package)" ]; then \
		echo "Usage: make test-package package=<package_path>"; \
		exit 1; \
	fi
	@mkdir -p tmp
	go test -v -coverprofile=tmp/coverage.out $(package)
	go tool cover -html=tmp/coverage.out -o tmp/coverage.html
	@echo "Coverage report generated at tmp/coverage.html"
	@open tmp/coverage.html

# Menjalankan test untuk domain user dengan coverage
test-user:
	@mkdir -p tmp
	go test -v -coverprofile=tmp/coverage.out ./internal/user/model/...
	go tool cover -html=tmp/coverage.out -o tmp/coverage.html
	@echo "Coverage report generated at tmp/coverage.html"
	@open tmp/coverage.html

# Membuka file coverage.html di browser
open-coverage:
	@if [ ! -f "tmp/coverage.html" ]; then \
		echo "Coverage report not found. Please run test-coverage first."; \
		exit 1; \
	fi
	@open tmp/coverage.html 