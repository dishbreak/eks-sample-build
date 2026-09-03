# Client: store_backend (Resource Server / API)
resource "keycloak_openid_client" "store_backend" {
  realm_id                     = keycloak_realm.contoso.id
  client_id                    = "store_backend"
  enabled                      = true
  access_type                  = "CONFIDENTIAL"
  standard_flow_enabled        = false
  direct_access_grants_enabled = false
  service_accounts_enabled     = true
  valid_redirect_uris          = ["/*"]
  use_refresh_tokens           = false
}

# Client: store_frontend (Consumer Service with store.read scope)
resource "keycloak_openid_client" "store_frontend" {
  realm_id                     = keycloak_realm.contoso.id
  client_id                    = "store_frontend"
  name                         = "Store Frontend Service"
  enabled                      = true
  always_display_in_console    = true
  access_type                  = "CONFIDENTIAL"
  standard_flow_enabled        = false
  direct_access_grants_enabled = false
  service_accounts_enabled     = true
  valid_redirect_uris          = ["/*"]
  use_refresh_tokens           = false
}

resource "keycloak_openid_client_default_scopes" "store_frontend_scopes" {
  realm_id  = keycloak_realm.contoso.id
  client_id = keycloak_openid_client.store_frontend.id

  default_scopes = [
    "web-origins",
    "acr",
    "profile",
    "roles",
    "email",
    keycloak_openid_client_scope.store_read.name,
    keycloak_openid_client_scope.store_backend_service.name,
  ]
}

# Client: store_keeper (Admin Consumer Service with store.read and store.admin scopes)
resource "keycloak_openid_client" "store_keeper" {
  realm_id                     = keycloak_realm.contoso.id
  client_id                    = "store_keeper"
  enabled                      = true
  always_display_in_console    = true
  access_type                  = "CONFIDENTIAL"
  standard_flow_enabled        = false
  direct_access_grants_enabled = false
  service_accounts_enabled     = true
  valid_redirect_uris          = ["/*"]
  use_refresh_tokens           = false
}

resource "keycloak_openid_client_default_scopes" "store_keeper_scopes" {
  realm_id  = keycloak_realm.contoso.id
  client_id = keycloak_openid_client.store_keeper.id

  default_scopes = [
    "web-origins",
    "acr",
    "profile",
    "roles",
    "email",
    keycloak_openid_client_scope.store_read.name,
    keycloak_openid_client_scope.store_admin.name,
    keycloak_openid_client_scope.store_backend_service.name,
  ]
}
