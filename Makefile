.PHONY: build run

build:
	GOOS=js GOARCH=wasm go build -o web/main.wasm ./cmd/app

run: build
	go run ./cmd/listen
