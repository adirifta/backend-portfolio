#!/bin/bash

# Script untuk run database migration ke Supabase

set -e

echo "🚀 Starting Supabase database migration..."

# Load environment variables
if [ -f /app/.env.local ]; then
    export $(grep -v '^#' /app/.env.local | xargs)
else
    echo "⚠️  Warning: .env.local file not found, using environment variables"
fi

# Check required environment variables
if [ -z "$DB_HOST" ] || [ -z "$DB_USER" ] || [ -z "$DB_PASSWORD" ] || [ -z "$DB_NAME" ]; then
    echo "❌ Error: Database connection variables not set"
    echo "DB_HOST: $DB_HOST"
    echo "DB_USER: $DB_USER" 
    echo "DB_NAME: $DB_NAME"
    exit 1
fi

# Set PGPASSWORD environment variable (important!)
export PGPASSWORD=$DB_PASSWORD

# Wait for database to be ready
echo "⏳ Waiting for Supabase database to be ready..."
until pg_isready -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t 10; do
    echo "Waiting for database connection..."
    sleep 2
done

echo "✅ Database is ready!"

# Test connection first
echo "🔌 Testing database connection..."
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SELECT version();" -t

# Run migration scripts
if [ -f "/app/scripts/init.sql" ]; then
    echo "📦 Running initial database setup..."
    psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f /app/scripts/init.sql -v ON_ERROR_STOP=1
    echo "✅ Initial setup completed"
else
    echo "⚠️  init.sql not found, skipping initial setup"
fi

# Run any additional migrations
if [ -d "/app/scripts/migrations" ]; then
    echo "🔄 Running database migrations..."
    for migration_file in /app/scripts/migrations/*.sql; do
        if [ -f "$migration_file" ]; then
            echo "Running migration: $(basename $migration_file)"
            psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f "$migration_file" -v ON_ERROR_STOP=1
            echo "✅ Migration completed: $(basename $migration_file)"
        fi
    done
else
    echo "⚠️  No migrations directory found"
fi

echo "🎉 Database migration completed successfully!"