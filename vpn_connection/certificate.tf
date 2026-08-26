resource "tls_private_key" "this" {
  for_each  = toset(["ca", "server", "client"])
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "ca" {
  private_key_pem = tls_private_key.this["ca"].private_key_pem

  subject {
    common_name  = "personal-vpn-ca"
    organization = "Personal VPN"
  }

  validity_period_hours = 12
  early_renewal_hours   = 6

  is_ca_certificate = true

  allowed_uses = [
    "cert_signing",
    "crl_signing",
    "digital_signature",
    "key_encipherment"
  ]
}

resource "tls_cert_request" "server" {
  private_key_pem = tls_private_key.this["server"].private_key_pem

  subject {
    common_name  = "client-vpn-server"
    organization = "Personal VPN"
  }

  dns_names = ["client-vpn-server"]
}

resource "tls_cert_request" "client" {
  private_key_pem = tls_private_key.this["client"].private_key_pem

  subject {
    common_name  = "laptop-client"
    organization = "Personal VPN"
  }
}

resource "tls_locally_signed_cert" "server" {
  cert_request_pem      = tls_cert_request.server.cert_request_pem
  ca_private_key_pem    = tls_private_key.this["ca"].private_key_pem
  ca_cert_pem           = tls_self_signed_cert.ca.cert_pem
  validity_period_hours = 6
  early_renewal_hours   = 3
  allowed_uses = [
    "digital_signature",
    "key_encipherment",
    "server_auth",
  ]
}

resource "tls_locally_signed_cert" "client" {
  cert_request_pem      = tls_cert_request.client.cert_request_pem
  ca_private_key_pem    = tls_private_key.this["ca"].private_key_pem
  ca_cert_pem           = tls_self_signed_cert.ca.cert_pem
  validity_period_hours = 6
  early_renewal_hours   = 3
  allowed_uses = [
    "digital_signature",
    "key_encipherment",
    "client_auth",
  ]
}

resource "aws_acm_certificate" "server" {
  private_key       = tls_private_key.this["server"].private_key_pem
  certificate_body  = tls_locally_signed_cert.server.cert_pem
  certificate_chain = tls_self_signed_cert.ca.cert_pem
  tags = {
    Name = "personal-client-vpn-server"
  }
}

resource "aws_acm_certificate" "client" {
  private_key       = tls_private_key.this["client"].private_key_pem
  certificate_body  = tls_locally_signed_cert.client.cert_pem
  certificate_chain = tls_self_signed_cert.ca.cert_pem
  tags = {
    name = "personal-client-vpn-laptop"
  }
}

resource "aws_acm_certificate" "client_ca" {
  private_key      = tls_private_key.this["ca"].private_key_pem
  certificate_body = tls_self_signed_cert.ca.cert_pem
  tags = {
    Name = "personal-client-vpn-ca"
  }
}

locals {
  client_private_key_pem = tls_private_key.this["client"].private_key_pem
  client_signed_cert_pem = tls_locally_signed_cert.client.cert_pem
}

output "client_private_key_pem" {
  value     = tls_private_key.this["client"].private_key_pem
  sensitive = true
}

output "client_signed_cert_pem" {
  value     = tls_locally_signed_cert.client.cert_pem
  sensitive = true
}
