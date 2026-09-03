output "issuer_url" {
  description = "The OIDC discovery issuer URL for this realm"
  value       = "${var.keycloak_url}/realms/${keycloak_realm.contoso.realm}"
}

output "store_backend_client_id" {
  description = "Client ID for the backend API"
  value       = keycloak_openid_client.store_backend.client_id
}

output "store_backend_client_secret" {
  description = "Client secret for store_backend"
  value       = keycloak_openid_client.store_backend.client_secret
  sensitive   = true
}

output "store_frontend_client_id" {
  description = "Client ID for store_frontend service"
  value       = keycloak_openid_client.store_frontend.client_id
}

output "store_frontend_client_secret" {
  description = "Client secret for store_frontend"
  value       = keycloak_openid_client.store_frontend.client_secret
  sensitive   = true
}

output "store_keeper_client_id" {
  description = "Client ID for store_keeper service"
  value       = keycloak_openid_client.store_keeper.client_id
}

output "store_keeper_client_secret" {
  description = "Client secret for store_keeper"
  value       = keycloak_openid_client.store_keeper.client_secret
  sensitive   = true
}
