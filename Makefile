.PHONY: build run cert

build:
	GOOS=js GOARCH=wasm go build -o web/main.wasm ./cmd/app

cert: cert.pem

cert.pem key.pem:
	openssl req -x509 -newkey ec -pkeyopt ec_paramgen_curve:prime256v1 -nodes \
		-days 3650 -subj "/CN=neringa" \
		-addext "subjectAltName=DNS:localhost,IP:127.0.0.1" \
		-keyout key.pem -out cert.pem

run: build cert.pem
	go run ./cmd/listen
