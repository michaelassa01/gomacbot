# Simple Makefile for a Go project
SHELL := /bin/bash
PROJECTNAME := "gomacbot"


# Build the application
all: build test

build:
	@echo "Building..."
	
	
	@go build -o main cmd/api/main.go

install-dependencies:
# 	sudo apt install ca-certificates
#gomigrate
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/swaggo/swag/cmd/swag@latest
#for Gin
	go get github.com/swaggo/gin-swagger
	go get github.com/swaggo/files
#firebase oauth
	go get firebase.google.com/go/v4
	go get google.golang.org/api/option

install-dependencies-wins:
	sudo apt install ca-certificates


install-dependencies-linux:
	sudo apt update
	sudo apt install -y ca-certificates
	sudo update-ca-certificates

# Run the application
run:
	@go run cmd/api/main.go
# Create DB container
docker-run:
	@if docker compose up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v
# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

# Usage: make module name=users
# module:
# 	@if [ -z "$(name)" ]; then \
# 		echo "Please provide a module name. Example: make module name=users"; \
# 		exit 1; \
# 	fi; \
# 	echo "📦 Creating module: $(name)"; \
# 	BASE_DIR=internal/$(name); \
# 	VERSION=$$(printf "%06d" $$(ls -1 $$BASE_DIR/migration 2>/dev/null | wc -l)); \
#     MIGRATION_NAME="$${VERSION}_init_schema"; \
# 	mkdir -p $$BASE_DIR/{migration,query,models}; \
# 	for f in router.go repository.go service.go handler.go; do \
# 		echo "package $(name)" > $$BASE_DIR/$$f; \
# 	done; \
# 	migrate create -ext sql -dir $$BASE_DIR/migration $$MIGRATION_NAME; \
# 	echo "✅ Module structure created for $(name)"; \
# 	\
# 	# ✅ Ensure sqlc.yaml exists and initialized \
# 	if [ ! -f ./sqlc.yaml ] || [ ! -s ./sqlc.yaml ]; then \
# 		printf "version: \"2\"\ncloud:\n  project: \"$(PROJECTNAME)\"\nsql:\n" > ./sqlc.yaml; \
# 		echo "✅ Created base sqlc.yaml"; \
# 	fi; \
# 	\
# 	# ✅ Append module block only if not already present \
# 	# Append module only if not already present
# 	if ! grep -q "schema: ./internal/$(name)/migration" ./sqlc.yaml; then
# 		cat >> ./sqlc.yaml <<EOL
# 	- engine: "postgresql"
# 		schema: "./internal/$(name)/migration"
# 		queries: "./internal/$(name)/query"

# 		gen:
# 		go:
# 			package: "$(name)"
# 			out: "./internal/$(name)/models"
# 			sql_package: "pgx/v5"
# 			emit_json_tags: true
# 			emit_prepared_queries: false
# 			emit_interface: true
# 			emit_exact_table_names: false
# 			emit_empty_slices: true
# 			overrides:
# 			- db_type: timestamptz
# 				go_type: time.Time

# 	EOL
# 		echo "✅ Added SQLC config for module $(name)"
# 	else
# 		echo "⚠️ SQLC config for module $(name) already exists — skipping"
# 	fi


module:
	@if [ -z "$(name)" ]; then \
		echo "Please provide a module name. Example: make module name=users"; \
		exit 1; \
	fi
	@bash scripts/create_sqlc_module.sh "$(name)" "$(PROJECTNAME)"





# to tidy the code base by installing dependencies
tidy:
	go mod tidy

migration:
	@if [ -z "$(module)" ]; then \
		echo "❌ Please provide a module. Example: make migration module=users name=add_profile_table"; \
		exit 1; \
	fi
	@if [ -z "$(name)" ]; then \
		echo "❌ Please provide a migration name. Example: make migration module=users name=add_profile_table"; \
		exit 1; \
	fi
	@if [ ! -d "internal/$(module)" ]; then \
		echo "❌ Module '$(module)' does not exist in internal/"; \
		exit 1; \
	fi

	@BASE_DIR=internal/$(module); \
	MIGRATION_DIR=$$BASE_DIR/migration; \
	mkdir -p $$MIGRATION_DIR; \
	TIMESTAMP=$$(date +"%Y%m%d%H%M%S"); \
	COUNT=$$(ls -1 $$MIGRATION_DIR/*.up.sql 2>/dev/null | wc -l | tr -d ' '); \
	SEQ=$$(printf "%06d" $$((COUNT + 1))); \
	FILENAME="$${TIMESTAMP}_$${SEQ}_$(name)"; \
	echo "🚀 Creating migration $$FILENAME for module '$(module)'"; \
	touch $$MIGRATION_DIR/$$FILENAME.up.sql; \
	touch $$MIGRATION_DIR/$$FILENAME.down.sql; \
	echo "✅ Created:"; \
	echo "   - $$MIGRATION_DIR/$$FILENAME.up.sql"; \
	echo "   - $$MIGRATION_DIR/$$FILENAME.down.sql";

sqlc:
	sqlc generate

db_docs:
	dbdocs build internal/doc/db.dbml

db_schema:
	dbml2sql internal/doc/db.dbml -o internal/doc/db.sql

swagger-doc:
	swag init -g cmd/main.go

mock:
	mockgen -destination db/mock/store.go  github.com/michaelassa01/gomacbot/db/models Store

.PHONY: all build run test clean watch docker-run docker-down itest install-dependencies test tidy swagger-doc createmigration module sqlc migration db_docs db_schema mock
