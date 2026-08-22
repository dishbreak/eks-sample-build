module "test_instances" {
  source     = "./module/test_instances"
  count      = var.provision_test_instances ? 1 : 0
  subnet_ids = { for k, v in aws_subnet.this : k => v.id }
  vpc_id     = aws_vpc.this.id
}
