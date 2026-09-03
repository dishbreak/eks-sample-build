# Keycloak Terraform Configuration

This directory contains Terraform configuration for provisioning the Keycloak resources used by the `store_backend` service using the [`mrparkers/keycloak`](https://registry.terraform.io/providers/mrparkers/keycloak/latest/docs) provider.

## Resources Managed

- **Realm:** `contoso`
- **Client Scopes:**
  - `store.read` (`include_in_token_scope = true`)
  - `store.admin` (`include_in_token_scope = true`)
  - `store_backend-service` (contains audience mapper adding `store_backend` to `aud`)
- **Clients:**
  - `store_backend`: Resource Server / API client
  - `store_frontend`: Service client with default scopes including `store.read` and `store_backend-service`
  - `store_keeper`: Admin service client with default scopes including `store.read`, `store.admin`, and `store_backend-service`

## Usage

```bash
# Initialize provider
terraform init

# Plan changes
terraform plan

# Apply to local Keycloak
terraform apply
```

To target a remote Keycloak instance or non-default credentials:

```bash
terraform apply \
  -var="keycloak_url=http://keycloak:7080" \
  -var="keycloak_admin_username=admin" \
  -var="keycloak_admin_password=admin"
```
