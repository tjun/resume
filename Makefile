.PHONY: build serve pdf lint clean

build:
	hugo

serve:
	hugo server

pdf: build
	go run ./cmd/pdf

lint:
	yarn lint

clean:
	rm -rf public
	rm -f resume.pdf
