#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

MODULE_NAME="$1"
PROJECT_NAME="$2"

# Check if module name is provided
if [ -z "$MODULE_NAME" ]; then
    echo -e "${RED}ERROR: Module name is required${NC}"
    echo "Usage: $0 <module_name> [project_name]"
    echo "Example: $0 users"
    exit 1
fi

BASE_DIR="internal/$MODULE_NAME"

# Check if module already exists
if [ -d "$BASE_DIR" ]; then
    echo -e "${RED}ERROR: Module '$MODULE_NAME' already exists at $BASE_DIR${NC}"
    echo "Choose a different module name or delete the existing module first"
    echo "To delete the existing module: ./scripts/delete_module.sh $MODULE_NAME"
    exit 1
fi

echo -e "${GREEN}Creating module: $MODULE_NAME${NC}"

# Create directories
mkdir -p "$BASE_DIR/migration" "$BASE_DIR/query" "$BASE_DIR/models" "$BASE_DIR/domain"
echo -e "${GREEN}Directories created${NC}"

# Create Go files
for f in router.go repository.go service.go handler.go pg_repo.go dto.go; do
  echo "package $MODULE_NAME" > "$BASE_DIR/$f"
done
echo -e "${GREEN}Go files created${NC}"

# Create initial migration (assuming migrate CLI installed)
VERSION=$(printf "%06d" $(ls -1 "$BASE_DIR/migration" 2>/dev/null | wc -l))
MIGRATION_NAME="${VERSION}_init_schema"
if command -v migrate &> /dev/null; then
    migrate create -ext sql -dir "$BASE_DIR/migration" "$MIGRATION_NAME"
    echo -e "${GREEN}Migration created: $MIGRATION_NAME${NC}"
else
    echo -e "${YELLOW}Warning: 'migrate' command not found. Skipping migration creation${NC}"
fi

# Ensure sqlc.yaml exists
if [ ! -f sqlc.yaml ] || [ ! -s sqlc.yaml ]; then
cat > sqlc.yaml <<EOL
version: "2"
cloud:
  project: "$PROJECT_NAME"
sql:
EOL
  echo -e "${GREEN}Created base sqlc.yaml${NC}"
fi

# Append module to sqlc.yaml if not already present
if ! grep -q "schema: ./internal/$MODULE_NAME/migration" sqlc.yaml; then
cat >> sqlc.yaml <<EOL
  - engine: "postgresql"
    schema: "./internal/$MODULE_NAME/migration"
    queries: "./internal/$MODULE_NAME/query"

    gen:
      go:
        package: "$MODULE_NAME"
        out: "./internal/$MODULE_NAME/models"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_prepared_queries: false
        emit_interface: true
        emit_exact_table_names: false
        emit_empty_slices: true
        emit_pointers_for_null_types: true
        overrides:
          - db_type: timestamptz
            go_type: time.Time
EOL
  echo -e "${GREEN}Added SQLC config for module $MODULE_NAME${NC}"
else
  echo -e "${YELLOW}SQLC config for module $MODULE_NAME already exists — skipping${NC}"
fi

echo ""
echo -e "${GREEN}Module '$MODULE_NAME' created successfully!${NC}"
echo ""
echo "Created structure:"
echo "  $BASE_DIR/"
echo "  ├── migration/     # Database migrations"
echo "  ├── query/         # SQLC query files"
echo "  ├── models/        # Generated models"
echo "  ├── domain/        # Domain logic"
echo "  ├── router.go      # Route definitions"
echo "  ├── repository.go  # Repository interface"
echo "  ├── service.go     # Business logic"
echo "  ├── handler.go     # HTTP handlers"
echo "  ├── pg_repo.go     # PostgreSQL implementation"
echo "  └── dto.go         # Data transfer objects"
echo ""
echo "Next steps:"
echo "  1. Write your queries in $BASE_DIR/query/"
echo "  2. Run 'make sqlc' to generate models"
echo "  3. Implement your business logic in service.go"
echo "  4. Register routes in router.go"
