locals {
  vpc_cidr = "10.12.0.0/16"
  vpn_cidr = "10.13.0.0/16"
  availability_zones = [
    "us-west-2a", "us-west-2b", "us-west-2c",
  ]
  public_network_names = [
    for z in local.availability_zones : "${z}-public"
  ]
  private_network_names = [
    for z in local.availability_zones : "${z}-private"
  ]
  networks = flatten([
    [for n in local.public_network_names : { name = n, new_bits = 8 }],
    [for n in local.private_network_names : { name = n, new_bits = 8 }],
  ])
}

module "subnets" {
  source          = "hashicorp/subnets/cidr"
  base_cidr_block = local.vpc_cidr
  networks        = local.networks
}
