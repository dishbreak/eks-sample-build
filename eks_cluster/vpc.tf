locals {
  vpc_cidr = "10.10.0.0/16"
  availability_zones = [
    "us-west-2a",
    "us-west-2b",
    "us-west-2c",
    "us-west-2d",
  ]
  network_names = flatten([for z in local.availability_zones: [for s in ["public", "private"]: join("-", [z, s])]])
}

module "subnet_addrs" {
  source = "hashicorp/subnets/cidr"
  base_cidr_block = local.vpc_cidr
  networks = [for n in local.network_names : {name = n, new_bits = 8}]
}

resource "aws_vpc" "this" {
  cidr_block = local.vpc_cidr
  tags = {
    Name = "eks_auto_mode"
  }
}

resource "aws_subnet" "this" {
  for_each = module.subnet_addrs.network_cidr_blocks
  vpc_id = aws_vpc.this.id
  availability_zone = regex("^us-west-2[a-d]", each.key)
  cidr_block = each.value
  # if the network name has a public suffix, it's a public subnet
  map_public_ip_on_launch = endswith(each.key, "-public")
}

locals {
  private_subnets = {for k, v in aws_subnet.this : k => v if endswith(k, "-private")}
  public_subnets = {for k, v in aws_subnet.this : k => v if endswith(k, "-public")}
}

resource "aws_internet_gateway" "this" {
  vpc_id = aws_vpc.this.id
  tags = {
    Name = "eks_auto_mode"
  }
}

resource "aws_eip" "nat_gateway" {
  for_each = local.public_subnets
  domain = "vpc"
}

resource "aws_nat_gateway" "this" {
  for_each = local.public_subnets
  subnet_id = each.value.id
  tags = {
    Name = "eks_auto_mode"
  }
  allocation_id = aws_eip.nat_gateway[each.key].allocation_id
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.this.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.this.id
  }
}

resource "aws_route_table" "private" {
  for_each = local.private_subnets
  vpc_id = aws_vpc.this.id
  route {
    cidr_block = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.this[
      join("-", [regex("^us-west-2[a-d]", each.key), "public"])
    ].id
  }
}

resource "aws_route_table_association" "public" {
  for_each = local.public_subnets
  route_table_id = aws_route_table.public.id
  subnet_id = each.value.id
}

resource "aws_route_table_association" "private" {
  for_each = local.private_subnets
  route_table_id = aws_route_table.private[each.key].id
  subnet_id = each.value.id
}
