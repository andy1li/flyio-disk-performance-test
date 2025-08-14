.PHONY: help deploy logs run-local clean

# Default target
help:
	@echo "Available commands:"
	@echo "  deploy      - Deploy the application to Fly.io"
	@echo "  logs        - View logs from the deployed application"
	@echo "  run-local   - Run the application locally"
	@echo "  clean       - Clean up local build artifacts"
	@echo "  status      - Check the status of the deployed application"
	@echo "  destroy     - Destroy the deployed application"

# Deploy to Fly.io
deploy:
	@echo "🚀 Deploying to Fly.io..."
	fly deploy

# View logs from the deployed application
logs:
	@echo "📋 Viewing logs from Fly.io..."
	fly logs

# Follow logs in real-time
logs-follow:
	@echo "📋 Following logs from Fly.io..."
	fly logs --follow

# Run the application locally
run-local:
	@echo "🏃 Running application locally..."
	go run main.go

# Check the status of the deployed application
status:
	@echo "📊 Checking application status..."
	fly status

# Destroy the deployed application
destroy:
	@echo "🗑️  Destroying application on Fly.io..."
	fly destroy
