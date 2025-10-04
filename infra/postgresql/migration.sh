#!/bin/bash

set -euo pipefail

PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null)"
MIGRATIONS_DIR="$PROJECT_ROOT/infra/postgresql/migrations"
DB_URI="${DB_URI:-}"

echo "PROJECT_ROOT: $PROJECT_ROOT"
echo "MIGRATIONS_DIR: $MIGRATIONS_DIR"

if [[ -z "$DB_URI" ]]; then
    echo "Error: DB_URI environment variable not set"
    exit 1
fi

if [[ ! -d "$MIGRATIONS_DIR" ]]; then
    echo "Error: Migrations directory $MIGRATIONS_DIR not found"
    exit 1
fi

echo "🔍 Checking SQL files for idempotency issues..."

# Simple idempotency checks
check_idempotency() {
    local issues=0
    
    # Check CREATE TABLE without IF NOT EXISTS
    if grep -rn "CREATE TABLE " "$MIGRATIONS_DIR" | grep -v "IF NOT EXISTS"; then
        echo "⚠️  Found CREATE TABLE without IF NOT EXISTS"
        issues=$((issues + 1))
    fi
    
    # Check CREATE INDEX without IF NOT EXISTS  
    if grep -rn "CREATE.*INDEX " "$MIGRATIONS_DIR" | grep -v "IF NOT EXISTS"; then
        echo "⚠️  Found CREATE INDEX without IF NOT EXISTS"
        issues=$((issues + 1))
    fi
    
    # Check INSERT without conflict handling
    if grep -rn "INSERT INTO" "$MIGRATIONS_DIR" | grep -v -E "(ON CONFLICT|ON DUPLICATE)"; then
        echo "⚠️  Found INSERT without conflict handling"
        issues=$((issues + 1))
    fi
    
    # Check ALTER TABLE ADD COLUMN without IF NOT EXISTS (Postgres 9.6+)
    if grep -rn "ALTER TABLE.*ADD COLUMN" "$MIGRATIONS_DIR" | grep -v "IF NOT EXISTS"; then
        echo "⚠️  Found ALTER TABLE ADD COLUMN without IF NOT EXISTS"
        issues=$((issues + 1))
    fi
    
    # Check DROP without IF EXISTS
    if grep -rn "DROP TABLE\|DROP INDEX\|DROP CONSTRAINT" "$MIGRATIONS_DIR" | grep -v "IF EXISTS"; then
        echo "⚠️  Found DROP statement without IF EXISTS"
        issues=$((issues + 1))
    fi
    
    return $issues
}

# Run idempotency checks
if ! check_idempotency; then
    echo "❌ Found potential idempotency issues (warnings only)"
fi

echo "🔧 Checking Atlas availability..."

# Check if atlas is available
if ! command -v atlas &> /dev/null; then
    echo "Error: atlas not found. Install with: curl -sSf https://atlasgo.sh | sh"
    exit 1
fi

# Check if database is reachable
echo "🔗 Checking database connection..."
if ! atlas schema inspect --url "$DB_URI" > /dev/null 2>&1; then
    echo "❌ Cannot connect to database at $DB_URI"
    exit 1
fi

echo "📋 Checking for missing migration files..."

# Check for missing migration files
check_missing_migrations() {
    # Get applied migrations from database
    local applied_migrations
    applied_migrations=$(atlas migrate status --url "$DB_URI" --dir "file://$MIGRATIONS_DIR" 2>/dev/null | grep "Applied" | awk '{print $1}' || true)
    
    if [[ -n "$applied_migrations" ]]; then
        local missing_count=0
        while IFS= read -r migration_version; do
            # Look for migration files matching this version
            if ! find "$MIGRATIONS_DIR" -name "*${migration_version}*" -type f | grep -q .; then
                echo "⚠️  Missing migration file for applied version: $migration_version"
                missing_count=$((missing_count + 1))
            fi
        done <<< "$applied_migrations"
        
        if [[ $missing_count -gt 0 ]]; then
            echo "🚨 Found $missing_count missing migration files!"
            echo "   These migrations were applied to the database but files are missing."
            echo "   Database changes from these migrations are still active."
            echo "   Consider creating 'undo' migrations if you want to revert changes."
        fi
    fi
}

# Run missing migration check
check_missing_migrations

echo "📋 Ensuring migration metadata is up to date..."

# Run atlas migrate hash to ensure migration files have proper checksums
# This is idempotent - won't change anything if hashes already exist
if ! atlas migrate hash --dir "file://$MIGRATIONS_DIR"; then
    echo "❌ Failed to generate migration hashes"
    exit 1
fi

echo "🚀 Applying migrations..."

# Apply migrations (this is idempotent - won't reapply already applied migrations)
if atlas migrate apply \
    --url "$DB_URI" \
    --dir "file://$MIGRATIONS_DIR"; then
    echo "✅ Migrations completed successfully"
else
    echo "❌ Migration failed"
    exit 1
fi

# Optional: Show current migration status
echo "📊 Current migration status:"
atlas migrate status --url "$DB_URI" --dir "file://$MIGRATIONS_DIR" || echo "Could not retrieve status"

