test:
	go test -C pkg -v

build:
	go build -C cmd -o ../bin/gobj

start-dev-deps:
	docker compose -f compose/docker-compose.yaml up -d

run-api:
	go run cmd/tui.go --apiserver

run-tui:
	./bin/gobj
