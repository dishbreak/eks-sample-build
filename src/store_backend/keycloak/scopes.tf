# Client Scope: store.read
resource "keycloak_openid_client_scope" "store_read" {
  realm_id               = keycloak_realm.contoso.id
  name                   = "store.read"
  include_in_token_scope = true
}

# Client Scope: store.admin
resource "keycloak_openid_client_scope" "store_admin" {
  realm_id               = keycloak_realm.contoso.id
  name                   = "store.admin"
  include_in_token_scope = true
}

# Client Scope: store_backend-service (Audience Scope)
resource "keycloak_openid_client_scope" "store_backend_service" {
  realm_id               = keycloak_realm.contoso.id
  name                   = "store_backend-service"
  include_in_token_scope = false
}

# Audience Mapper: adds "store_backend" to the token "aud" claim
resource "keycloak_openid_audience_protocol_mapper" "store_backend_audience" {
  realm_id                 = keycloak_realm.contoso.id
  client_scope_id          = keycloak_openid_client_scope.store_backend_service.id
  name                     = "store_backend-audience"
  included_client_audience = "store_backend"
  add_to_access_token      = true
  add_to_id_token          = false
}
