#!/bin/bash

YAML_FILE="issues.yaml"

# Check for required tools
for cmd in gh jq yq; do
  if ! command -v $cmd &> /dev/null; then
    echo "❌ $cmd is not installed. Please install it before running this script."
    exit 1
  fi
done

# Read each issue from YAML and create it via GitHub CLI
yq e '.issues[]' -o=json "$YAML_FILE" | jq -c '.' | while read -r issue; do
  title=$(echo "$issue" | jq -r '.title')
  body=$(echo "$issue" | jq -r '.body')
  labels=$(echo "$issue" | jq -r '.labels | join(",")')
  milestone=$(echo "$issue" | jq -r '.milestone')

  echo "📦 Creating issue: $title"
  gh issue create --title "$title" --body "$body" --label "$labels" --milestone "$milestone"
done