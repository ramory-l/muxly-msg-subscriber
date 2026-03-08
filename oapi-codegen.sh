#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SCHEMAS_DIR="$SCRIPT_DIR/schemas"

echo "Generating Go code for muxly msg subscriber..."

echo "  - Generating msg subscriber server..."
oapi-codegen --config "$SCHEMAS_DIR/codegen/msg-subscriber-server.yaml" "$SCHEMAS_DIR/internal/muxly_msg_subscriber/api.yaml"

echo "  - Generating msg subscriber types..."
oapi-codegen --config "$SCHEMAS_DIR/codegen/msg-subscriber-types.yaml" "$SCHEMAS_DIR/internal/muxly_msg_subscriber/definitions.yaml"


echo "  - Generating muxly-core types..."
oapi-codegen --config "$SCHEMAS_DIR/codegen/muxly-core-types.yaml" "$SCHEMAS_DIR/internal/core_definitions.yaml"

echo "  - Generating muxly-backend client..."
oapi-codegen --config "$SCHEMAS_DIR/codegen/muxly-backend-client-for-msg-subscriber.yaml" "$SCHEMAS_DIR/internal/muxly_backend/api.yaml"

echo "  - Generating muxly-backend types..."
oapi-codegen --config "$SCHEMAS_DIR/codegen/muxly-backend-types.yaml" "$SCHEMAS_DIR/internal/muxly_backend/definitions.yaml"

echo "muxly msg subscriber code generation complete!"
