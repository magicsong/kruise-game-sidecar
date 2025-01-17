#!/bin/bash

# Usage: ./deploy.sh <new_image>

if [ $# -ne 1 ]; then
    echo "Usage: $0 <new_image>"
    exit 1
fi

NEW_IMAGE="$1"

FILE="config/manager/env_path.yaml"

sed -i '' "s|\\(value:[[:space:]]*\\).*|\\1$NEW_IMAGE|" "$FILE"


