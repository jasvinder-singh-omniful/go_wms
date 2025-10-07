APP_NAME := oms
BIN_PATH := ./tmp/main
SRC_PATH := ./cmd/server
CONFIG := config/local.yml
MASTER_DB_URL := postgres://postgres:123@localhost:5432/wms?sslmode=disable

# test makefile formatting of tabs (^I denotes tab and $ denotes line ending)
check:
	cat -e -t -v Makefile

# build go binary 
build:
	go build -o $(BIN_PATH) $(SRC_PATH)


# migrations commands
create-migration:
	migrate create -ext sql -dir migrations -seq create_tables

migrate-up:
	@echo "applying migrations to MASTER database..."
	migrate -database $(MASTER_DB_URL) -path migrations up

migrate-down:
	@echo "rolling back migrations from MASTER database..."
	migrate -database $(MASTER_DB_URL) -path migrations down

migrate-version:
	@echo "current migration version on MASTER database..."
	migrate -database $(MASTER_DB_URL) -path migrations version

migrate-force:
	@migrate -database "$(MASTER_DB_URL)" -path migrations force $(version)

# run the dev server
run-dev:
	docker stop postgres || true
	docker stop redis || true
	docker compose -f docker/docker-compose.yml up -d
	air

# clean up binary
clean:
	rm -f $(BIN_PATH)

# run tests
test:
	go test ./...

.PHONY: build run clean test check