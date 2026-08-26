locals {
  subnet_ids = merge({ for k, v in aws_subnet.public : k => v.id }, { for k, v in aws_subnet.private : k => v.id })
}

module "test_instances" {
  source     = "../eks_cluster/module/test_instances"
  vpc_id     = aws_vpc.this.id
  subnet_ids = local.subnet_ids
}
