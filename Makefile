
# Postgres docker container is based on port 5433
migrate-up:
	migrate -path ./dbschema -database postgres://postgres:qwerty123@0.0.0.0:5433/postgres?sslmode=disable up

migrate-down:
	migrate -path ./dbschema -database postgres://postgres:qwerty123@0.0.0.0:5433/postgres?sslmode=disable down 1

migrate-drop:
	migrate -path ./dbschema -database postgres://postgres:qwerty123@0.0.0.0:5433/postgres?sslmode=disable drop

lint:
	golangci-lint --config .golangci.yml run ./... --deadline=2m --timeout=2m
