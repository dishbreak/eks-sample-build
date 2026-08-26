resource "aws_security_group" "vpn_access" {
  name        = "vpn-access"
  description = "security group for AWS VPN Client"
  vpc_id      = aws_vpc.this.id
}

resource "aws_security_group_rule" "vpn_outbound_all" {
  security_group_id = aws_security_group.vpn_access.id
  type              = "egress"
  from_port         = -1
  to_port           = -1
  protocol          = -1
  cidr_blocks       = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "icmp_to_test_inst" {
  security_group_id        = module.test_instances.security_group_id
  source_security_group_id = aws_security_group.vpn_access.id
  type                     = "ingress"
  from_port                = -1
  to_port                  = -1
  protocol                 = "icmp"
}

resource "aws_ec2_client_vpn_endpoint" "this" {
  vpc_id             = aws_vpc.this.id
  security_group_ids = [aws_security_group.vpn_access.id]
  authentication_options {
    type                       = "certificate-authentication"
    root_certificate_chain_arn = aws_acm_certificate.client_ca.arn
  }
  server_certificate_arn = aws_acm_certificate.server.arn
  connection_log_options {
    enabled = false
  }
  client_cidr_block = local.vpn_cidr
  split_tunnel      = true
}

resource "aws_ec2_client_vpn_network_association" "this" {
  for_each               = toset(local.private_network_names)
  subnet_id              = aws_subnet.private[each.key].id
  client_vpn_endpoint_id = aws_ec2_client_vpn_endpoint.this.id
}

resource "aws_ec2_client_vpn_authorization_rule" "full_vpc_access" {
  client_vpn_endpoint_id = aws_ec2_client_vpn_endpoint.this.id
  target_network_cidr    = local.vpc_cidr
  authorize_all_groups   = true
}

data "aws_region" "this" {}

resource "terraform_data" "vpn_config" {
  triggers_replace = [
    aws_ec2_client_vpn_endpoint.this.id,
    local.client_private_key_pem,
    local.client_signed_cert_pem,
  ]

  provisioner "local-exec" {
    environment = {
      PRIVATE_KEY_PEM    = local.client_private_key_pem
      CERT               = local.client_signed_cert_pem
      CLIENT_ENDPOINT_ID = aws_ec2_client_vpn_endpoint.this.id
      AWS_REGION         = data.aws_region.this.region
    }
    command     = "./generate_ovpn.sh"
    working_dir = "./scripts"
  }
}
