test:
	go test -C pkg -v

build:
	go build -C cmd -o ../bin/gobj

run-api:
	docker compose -f compose/docker-compose.yaml up &
	go run cmd/tui.go --apiserver

run-tui:
	./bin/gobj
