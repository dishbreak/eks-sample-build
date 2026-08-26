resource "aws_vpc" "this" {
  cidr_block = local.vpc_cidr
  tags = {
    "Name" = "vpn_connection_test"
  }
}

resource "aws_subnet" "public" {
  for_each                = toset(local.public_network_names)
  vpc_id                  = aws_vpc.this.id
  availability_zone       = trimsuffix(each.key, "-public")
  map_public_ip_on_launch = true
  cidr_block              = module.subnets.network_cidr_blocks[each.key]
  tags = {
    Name = each.key
  }
}

resource "aws_subnet" "private" {
  for_each                = toset(local.private_network_names)
  vpc_id                  = aws_vpc.this.id
  availability_zone       = trimsuffix(each.key, "-private")
  map_public_ip_on_launch = false
  cidr_block              = module.subnets.network_cidr_blocks[each.key]
  tags = {
    Name = each.key
  }
}

resource "aws_internet_gateway" "this" {
  vpc_id = aws_vpc.this.id
  tags = {
    Name = "vpn_connection_test"
  }
}

resource "aws_eip" "nat_gateway" {
  for_each = toset(local.public_network_names)
  domain   = "vpc"
  tags = {
    Name = "${each.key}-nat-gateway"
  }
}

resource "aws_nat_gateway" "name" {
  for_each      = toset(local.public_network_names)
  subnet_id     = aws_subnet.public[each.key].id
  allocation_id = aws_eip.nat_gateway[each.key].allocation_id
  tags = {
    Name = "${each.key}-nat-gateway"
  }
}

locals {
  zonal_gateways = { for k, v in aws_nat_gateway.name : trimsuffix(k, "-public") => v.id }
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.this.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.this.id
  }
}

resource "aws_route_table_association" "public" {
  for_each       = toset(local.public_network_names)
  route_table_id = aws_route_table.public.id
  subnet_id      = aws_subnet.public[each.key].id
}

resource "aws_route_table" "private" {
  for_each = toset(local.private_network_names)
  vpc_id   = aws_vpc.this.id
  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = local.zonal_gateways[trimsuffix(each.key, "-private")]
  }
}

resource "aws_route_table_association" "private" {
  for_each       = toset(local.private_network_names)
  route_table_id = aws_route_table.private[each.key].id
  subnet_id      = aws_subnet.private[each.key].id
}
