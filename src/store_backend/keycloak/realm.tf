resource "keycloak_realm" "contoso" {
  realm                       = var.realm_name
  enabled                     = true
  default_signature_algorithm = "RS256"
}
