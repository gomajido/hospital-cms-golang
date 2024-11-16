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
	mockgen --package mocks --source=pkg/storage/storage.go --destination=pkg/storage/mocks/mock_storage.go

	# driver
	@echo "  >  Driver Mocking..."
	mockgen  --package mocks --source=pkg/http_client/http/http_client.go --destination=pkg/http_client/http/mocks/mock_http_client.go

	# common
	@echo "  >  Common Mocking..."
	mockgen  --package mocks --source=internal/common/privy/domain/repository.go --destination=internal/common/privy/domain/mocks/mock_repository.go
	mockgen  --package mocks --source=internal/common/whatsapp/domain/repository.go --destination=internal/common/whatsapp/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/common/storage/domain/repository.go --destination=internal/common/storage/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/common/bsi_apexa/domain/repository.go --destination=internal/common/bsi_apexa/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/common/mailer/domain/repository.go --destination=internal/common/mailer/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/common/notification/domain/repository.go --destination=internal/common/notification/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/common/content/domain/repository.go --destination=internal/common/content/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/common/user/domain/repository.go --destination=internal/common/user/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/common/pdfconverter/domain/repository.go --destination=internal/common/pdfconverter/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/common/media/domain/repository.go --destination=internal/common/media/domain/mocks/mock_repository.go


	# modul
	@echo "  >  Usecase and Repository Mocking..."
	mockgen --package mocks --source=internal/module/auth/domain/repository.go --destination=internal/module/auth/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/auth/domain/usecase.go --destination=internal/module/auth/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/otp/domain/repository.go --destination=internal/module/otp/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/otp/domain/usecase.go --destination=internal/module/otp/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/register_individual/domain/usecase.go --destination=internal/module/register_individual/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/register_individual/domain/repository.go --destination=internal/module/register_individual/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/register/domain/repository.go --destination=internal/module/register/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/register/domain/usecase.go --destination=internal/module/register/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/register_issuer/domain/usecase.go --destination=internal/module/register_issuer/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/register_issuer/domain/repository.go --destination=internal/module/register_issuer/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/referral/domain/usecase.go --destination=internal/module/referral/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/referral/domain/repository.go --destination=internal/module/referral/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/user_setting/domain/repository.go --destination=internal/module/user_setting/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/user_setting/domain/usecase.go --destination=internal/module/user_setting/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/register_company/domain/usecase.go --destination=internal/module/register_company/domain/mocks/mock_usecase.go
	mockgen --package mocks --source=internal/module/register_company/domain/repository.go --destination=internal/module/register_company/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/user_setting/domain/repository.go --destination=internal/module/user_setting/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/user_setting/domain/usecase.go --destination=internal/module/user_setting/domain/mocks/mock_usecase.go

	mockgen --package mocks --source=internal/module/privy_registration/domain/repository.go --destination=internal/module/privy_registration/domain/mocks/mock_repository.go
	mockgen --package mocks --source=internal/module/privy_registration/domain/usecase.go --destination=internal/module/privy_registration/domain/mocks/mock_usecase.go



