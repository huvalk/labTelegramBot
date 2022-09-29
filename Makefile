build-app:
	CGO_CFLAGS="-g -O2 -Wno-return-local-addr" go build -o app cmd/main.go

start:
	docker-compose up -d
