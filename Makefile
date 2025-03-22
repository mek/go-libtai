#
LIBFILES=lib/tai/tai.go lib/taia/taia.go
TARGETS=tai_now taia_now tai64n
#
default:
	@echo "Targets"
	@echo " clean"
	@echo " $(TARGETS)"

.PHONY: clean all $(TARGETS) vet fmt lint check lib-test pre-commit
all: $(TARGETS)

check: fmt vet lint
	@pre-commit run --all-files

pre-commit:
	@pre-commit install
	@pre-commit autoupdate

lib-test:
	@go test -tags=unit ./lib/...

clean:
	@rm -f *~ bin/$(TARGETS) .*~

vet:
	@go vet ./...

fmt:
	@go fmt ./...

lint:
	@golangci-lint run ./...

tai_now: bin/tai_now
bin/tai_now: cmd/tai_now/main.go $(LIBFILES)
	@go build -o bin/tai_now ./cmd/tai_now

taia_now: bin/taia_now
bin/taia_now: cmd/taia_now/main.go $(LIBFILES)
	@go build -o bin/taia_now ./cmd/taia_now

tai64n: bin/tai64n
bin/tai64n: cmd/tai64n/main.go $(LIBFILES)
	@go build -o bin/tai64n ./cmd/tai64n
