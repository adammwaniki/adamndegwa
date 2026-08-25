#!/usr/bin/env bash
# Fails if any internal package drops below 100% test coverage.
# The root package (main.go) is an intentionally untestable 3-line shell.
set -euo pipefail

out=$(go test ./internal/... -cover)
echo "$out"

fail=0
while read -r line; do
  pct=$(echo "$line" | grep -oE '[0-9]+\.[0-9]+%' | head -1)
  pkg=$(echo "$line" | awk '{print $2}')
  if [[ "$pct" != "100.0%" ]]; then
    echo "FAIL: $pkg coverage $pct (want 100.0%)"
    fail=1
  fi
done < <(echo "$out" | grep 'coverage:')

if [[ $fail -ne 0 ]]; then
  exit 1
fi
echo "OK: all internal packages at 100% coverage"
