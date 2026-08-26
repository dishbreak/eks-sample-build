variable "vpc_id" {
  type        = string
  description = "the VPC to establish the VPN connection to"
}

variable "vpc_cidr" {
  type        = string
  description = "The CIDR corresponding to the target VPC"
}

variable "vpn_cidr" {
  type        = string
  description = "The CIDR corresponding to the addresses assigned to VPN clients. Must not overlap with VPC CIDR, and must be between /12 and /22."
}

variable "subnet_ids" {
  type        = map(string)
  description = "the subnet IDs to create associations with. Keys are friendly network names, values are subnet IDs."
}
