#!/bin/bash

# Load environment variables
if [ -f .env.local ]; then
    export $(cat .env.local | grep -v '#' | awk '/=/ {print $1}')
fi

# Check if required environment variables are set
if [ -z "$SUPABASE_PROJECT_REF" ] || [ -z "$SUPABASE_DB_PASSWORD" ]; then
    echo "Error: SUPABASE_PROJECT_REF and SUPABASE_DB_PASSWORD must be set in .env.local"
    exit 1
fi

# Build and run with Docker Compose
docker-compose -f docker-compose.supabase.yml up --build