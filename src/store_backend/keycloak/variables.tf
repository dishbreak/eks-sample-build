variable "keycloak_url" {
  description = "The URL of the Keycloak server"
  type        = string
  default     = "http://localhost:7080"
}

variable "keycloak_admin_username" {
  description = "Keycloak bootstrap admin username"
  type        = string
  default     = "admin"
}

variable "keycloak_admin_password" {
  description = "Keycloak bootstrap admin password"
  type        = string
  default     = "admin"
  sensitive   = true
}

variable "realm_name" {
  description = "The name of the Keycloak realm"
  type        = string
  default     = "contoso"
}
