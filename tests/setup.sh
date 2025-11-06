#!/bin/bash

# Moogie API Test Setup Script
echo "🧪 Setting up Moogie API Tests with Bruno"
echo "========================================"

# Check if Bruno is installed
if ! command -v bruno &> /dev/null; then
    echo "❌ Bruno is not installed."
    echo ""
    echo "Please install Bruno:"
    echo "  • macOS: brew install bruno"
    echo "  • Or download from: https://www.usebruno.com/"
    echo ""
    exit 1
fi

echo "✅ Bruno is installed"

# Check if API server is running
API_URL="http://localhost:8080/health"
echo "🔍 Checking if API server is running..."

if curl -s "$API_URL" > /dev/null 2>&1; then
    echo "✅ API server is running at http://localhost:8080"
else
    echo "❌ API server is not running"
    echo ""
    echo "Please start the API server:"
    echo "  cd api && make dev"
    echo "  # or"
    echo "  cd api && go run cmd/server/main.go"
    echo ""
    exit 1
fi

# Test the health endpoint
echo "🏥 Testing health endpoint..."
HEALTH_RESPONSE=$(curl -s "$API_URL")
if [[ $? -eq 0 ]]; then
    echo "✅ Health check successful: $HEALTH_RESPONSE"
else
    echo "❌ Health check failed"
    exit 1
fi

echo ""
echo "🎉 Setup complete! You can now:"
echo "   1. Open Bruno"
echo "   2. Click 'Open Collection'"
echo "   3. Navigate to $(pwd)/api"
echo "   4. Select 'local' environment"
echo "   5. Run your tests!"
echo ""
echo "💡 Tip: Check the README.md for detailed usage instructions"