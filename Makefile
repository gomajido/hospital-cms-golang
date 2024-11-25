APP_EXECUTABLE="out/apexa-service"

serve-api:
	export APEXA_ENV=local && \
	go run main.go serve-rest-api

compile:
	mkdir -p out/
	go build -o $(APP_EXECUTABLE)

test:
	go test -v ./...

test-coverage:
	./coverage.sh;

static-check:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./...

lint:
	@golangci-lint run -E gofmt

format:
	@$(MAKE) fmt
	@$(MAKE) imports

fmt:
	@echo "Formatting code style..."
	gofmt -w -s cmd/.. \
		internal/..

	@echo "[DONE] Formatting code style..."

imports:
	@echo "Formatting imports..."
	goimports -w -local go-core-apexa \
		cmd/.. \
		config/.. \
		internal/.. \
		pkg/..
	@echo "[DONE] Formatting imports..."

migrate-up:
	@echo "Running migration up..."
	goose --dir=gen/migrations postgres "user=$(user) password=$(password) dbname=$(dbname) host=$(host) port=$(port) sslmode=disable" up
	@echo $(user) $(password) $(dbname)
	@echo "[DONE] Running migration up.."

migrate-down:
	@echo "Running migration down..."
	goose --dir=gen/migrations postgres "user=$(user) password=$(password) dbname=$(dbname) host=$(host) port=$(port) sslmode=disable" down
	@echo $(user) $(password) $(dbname)
	@echo "[DONE] Running migration down.."

gen-mocks:
	@echo "  >  Rebuild Mocking..."
	# mockgen --package mocks --source=pkg/storage/storage.go --destination=pkg/storage/mocks/mock_storage.go

	# driver
	@echo "  >  Driver Mocking..."
	# mockgen  --package mocks --source=pkg/http_client/http/http_client.go --destination=pkg/http_client/http/mocks/mock_http_client.go

	# common
	@echo "  >  Common Mocking..."
	# mockgen  --package mocks --source=internal/common/privy/domain/repository.go --destination=internal/common/privy/domain/mocks/mock_repository.go
	


	# modul
	@echo "  >  Usecase and Repository Mocking..."
	# auth module
	mockgen --package mocks --source=internal/module/auth/domain/repository.go --destination=internal/module/auth/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/auth/domain/usecase.go --destination=internal/module/auth/domain/mocks/mock_usecase.go
	



