#!/bin/bash

# Simple development server script for Moogie API

echo "🚀 Starting Moogie API Development Server..."
echo ""
echo "   📡 Server: http://localhost:8080"
echo "   🔗 API: http://localhost:8080/api/v1"
echo "   ❤️  Health: http://localhost:8080/health"
echo "   🔌 WebSocket: ws://localhost:8080/ws"
echo ""

# Check if PostgreSQL is needed
if [ ! -f ".env" ]; then
    echo "⚠️  No .env file found. Creating from .env.example..."
    cp .env.example .env
fi

echo "🔧 Building server..."
go build -o bin/server cmd/server/main.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo ""
    echo "📚 Note: Make sure PostgreSQL is running on localhost:5432"
    echo "   - Database: moogie"
    echo "   - User: moogie"
    echo "   - Password: moogie"
    echo ""
    echo "🚀 Starting server... (Press Ctrl+C to stop)"
    echo ""
    ./bin/server
else
    echo "❌ Build failed!"
    exit 1
fi
