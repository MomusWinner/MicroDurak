# .PHONY: proto-gen-players-service
# proto-gen-players-service:
# 	@echo "Generate proto for players service"
# 	protoc \
# 		--proto_path=./internal/services/players/proto/ \
# 		--go_out=./internal/services/players/contracts/proto \
# 		--go-grpc_out=./internal/services/players/contracts/proto \
# 		./internal/services/players/proto/players.proto

.PHONY: sql
sql:
	sqlc generate

.PHONY: swagger
swagger:
	@echo "Generate swagger documentation for all services"
	@echo "Generating swagger for players service..."
	cd internal/services/players && swag init -g cmd/main.go -o delivery/http/docs --parseDependency --parseInternal
	@echo "Generating swagger for auth service..."
	cd internal/services/auth && swag init -g cmd/main.go -o delivery/http/docs --parseDependency --parseInternal
	@echo "Generating swagger for matchmaker service..."
	cd internal/services/matchmaker && swag init -g cmd/main.go -o delivery/http/docs --parseDependency --parseInternal
	@echo "Swagger documentation generated for all services"

.PHONY: test
test:
	go test -v -race -coverprofile=coverage.out ./...

.PHONY: cover
cover: test
	@echo "Generating HTML coverage report..."
	go tool cover -html=coverage.out -o coverage.html
