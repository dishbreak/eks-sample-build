variable "number_of_availability_zones" {
  type        = number
  description = "The number of availability zones to use, set -1 to use all of them. Note that each availability zone requires at least 1 Elastic IP."
}

variable "excluded_availability_zones" {
  type        = list(string)
  description = "the zones that should explicitly be excluded. Note that this may result in fewer AZs than local.number_of_availability_zones dictates."
  default     = []
}

variable "eks_vpc_cidr" {
  type        = string
  description = "The CIDR block for the VPC. Ensure the CIDR is large enough to accomodate private and public subnets for each availability zone"
  default     = "10.10.0.0/16"
}


variable "region" {
  type        = string
  description = "The region to use for the deployment"
  default     = "us-west-2"
}

variable "provision_test_instances" {
  type        = bool
  description = "create one test instance per subnet, useful for connectivity checking."
  default     = false
}

variable "cluster_name" {
  type        = string
  description = "the name to assign to the EKS cluster"
  default     = "eks-auto"
}

variable "eks_cluster_admin_users" {
  type        = list(string)
  description = "A list of principals that will gain cluster admin permissions"
}
