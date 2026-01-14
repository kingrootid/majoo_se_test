#!/bin/bash
# Script to install swag and generate API documentation

echo "Installing Swag CLI..."
go install github.com/swaggo/swag/cmd/swag@latest

echo "Generating API documentation..."
swag init -g cmd/api/main.go --output docs

if [ $? -eq 0 ]; then
    echo "API documentation generated successfully in docs/ directory"
else
    echo "Error generating documentation"
    exit 1
fi