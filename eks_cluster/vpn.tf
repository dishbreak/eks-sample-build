module "vpn_connection" {
  source     = "./module/vpn_connection"
  count      = var.provision_client_vpn_connection ? 1 : 0
  vpc_id     = aws_vpc.this.id
  vpc_cidr   = var.eks_vpc_cidr
  vpn_cidr   = var.vpn_cidr
  subnet_ids = { for k, v in local.private_subnets : k => v.id }
}
