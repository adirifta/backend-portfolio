#!/bin/bash
set -e

echo "🚀 Starting deployment..."

# Variables
APP_DIR="/home/$USER/portfolio-app"
APP_NAME="portfolio-app"
PORT=8080

# Navigate to app directory
cd $APP_DIR

# Stop existing application
echo "🛑 Stopping existing application..."
pkill -f $APP_NAME || true
sleep 3

# Backup current version
echo "📦 Creating backup..."
TIMESTAMP=$(date +%Y%m%d%H%M%S)
tar -czf ../portfolio-backup-$TIMESTAMP.tar.gz .

# Set environment variables
export $(grep -v '^#' .env.production | xargs)
export PORT=$PORT

# Start application
echo "🔧 Starting application..."
chmod +x $APP_NAME
nohup ./$APP_NAME > app.log 2>&1 &

# Wait for application to start
echo "⏳ Waiting for application to start..."
sleep 10

# Health check
if curl -f http://localhost:$PORT/health > /dev/null 2>&1; then
    echo "✅ Deployment successful!"
    echo "📊 Application logs:"
    tail -20 app.log
else
    echo "❌ Deployment failed!"
    echo "🔍 Checking logs..."
    cat app.log
    exit 1
fi