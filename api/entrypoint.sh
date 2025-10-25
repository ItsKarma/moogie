#!/bin/sh

echo "🔧 Starting Moogie API..."

# Wait for PostgreSQL to be ready
echo "⏳ Waiting for PostgreSQL to be ready..."
while ! pg_isready -h postgres -p 5432 -U moogie; do
  echo "   Waiting for PostgreSQL..."
  sleep 2
done
echo "✅ PostgreSQL is ready!"

# Run database migrations
echo "🗄️  Running database migrations..."
goose -dir /app/migrations postgres "postgres://moogie:moogie@postgres:5432/moogie?sslmode=disable" up

if [ $? -eq 0 ]; then
    echo "✅ Migrations completed successfully!"
else
    echo "❌ Migration failed!"
    exit 1
fi

# Start the API server
echo "🚀 Starting API server..."
./server
