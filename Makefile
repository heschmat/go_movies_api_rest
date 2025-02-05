run/api:
	go run ./cmd/api

check:
	curl -i localhost:4000/v1/healthcheck
