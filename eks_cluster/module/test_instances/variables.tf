variable "subnet_ids" {
  type        = map(string)
  description = "a map of subnet names and IDs to create test instances in"
}

variable "vpc_id" {
  type        = string
  description = "the VPC to launch the instances in"
}

output "security_group_id" {
  type        = string
  description = "the security group ID for the test instances group"
  value       = aws_security_group.inst.id
}
