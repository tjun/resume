.PHONY: build serve pdf lint clean

build:
	hugo -c docs

serve:
	hugo server -c docs

pdf: build
	go run ./cmd/pdf

lint:
	yarn lint

clean:
	rm -rf public
	rm -f resume.pdf
