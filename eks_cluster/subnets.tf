data "aws_availability_zones" "this" {
  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
  exclude_names = var.excluded_availability_zones
}

locals {
  vpc_cidr           = "10.10.0.0/16"
  availability_zones = slice(data.aws_availability_zones.this.names, 0, var.number_of_availability_zones)
  network_names      = flatten([for z in local.availability_zones : [for s in ["public", "private"] : join("-", [z, s])]])
}

module "subnet_addrs" {
  source          = "hashicorp/subnets/cidr"
  base_cidr_block = var.eks_vpc_cidr
  networks        = [for n in local.network_names : { name = n, new_bits = 8 }]
}
