#!/bin/bash
set -ex

mkdir -p binaries

for b in cmd/*/; do
    dir="./${b%*/}"
    file="./binaries/$(basename "$dir")"
    go build -o "$file" "$dir" 
done
