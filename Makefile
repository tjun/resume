.PHONY: build serve pdf fmt lint-go lint-text lint clean

build:
	hugo -c docs

serve:
	hugo server -c docs

pdf: build
	go run ./cmd/pdf

fmt:
	go fmt ./...

lint-go:
	go vet ./...
	go fmt -n ./...

lint-text:
	yarn lint

lint: lint-go lint-text

clean:
	rm -rf public
	rm -f resume.pdf
