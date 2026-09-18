#!/bin/bash

CLIENT_SECRET=$(jq -r '.outputs.store_frontend_client_secret.value' terraform.tfstate)

ACCESS_TOKEN=$(curl -s -X POST http://keycloak:7080/realms/contoso/protocol/openid-connect/token \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "grant_type=client_credentials" \
    -d "client_id=store_frontend" \
    -d "client_secret=${CLIENT_SECRET}" \
    -d "scope=store.read" | jq -r '.access_token')

echo " ### accessing items with auth"
curl -X GET http://localhost:8080/items \
    -H "Authorization: Bearer ${ACCESS_TOKEN}" -I

echo " ### trying to create an item"
curl -sS -D - -X POST http://localhost:8080/items \
    -H "Authorization: Bearer ${ACCESS_TOKEN}" \
    -H "Content-Type: application/json" \
    -d '{title: "foo", "description":"bar"}'
