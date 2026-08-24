resource "aws_eks_cluster" "this" {
  name = var.cluster_name
  access_config {
    authentication_mode = "API"
  }
  role_arn = aws_iam_role.aws_eks_cluster["cluster"].arn

  # auto mode settings
  compute_config {
    enabled       = true
    node_pools    = ["general-purpose"]
    node_role_arn = aws_iam_role.aws_eks_cluster["node"].arn
  }
  kubernetes_network_config {
    elastic_load_balancing {
      enabled = true
    }
  }
  storage_config {
    block_storage {
      enabled = true
    }
  }
  bootstrap_self_managed_addons = false

  vpc_config {
    endpoint_private_access = true
    endpoint_public_access  = true

    subnet_ids = [for s in local.private_subnets : s.id]
  }
  depends_on = [
    aws_iam_role_policy_attachment.cluster,
    aws_iam_role_policy_attachment.node,
  ]
}

data "aws_iam_policy_document" "aws_eks_node_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "aws_eks_cluster_assume_role" {
  statement {
    actions = ["sts:AssumeRole", "sts:TagSession"]
    principals {
      type        = "Service"
      identifiers = ["eks.amazonaws.com"]
    }
  }
}

locals {
  assume_role_policy = {
    cluster = data.aws_iam_policy_document.aws_eks_cluster_assume_role.json
    node    = data.aws_iam_policy_document.aws_eks_node_assume_role.json
  }
}

resource "aws_iam_role" "aws_eks_cluster" {
  for_each           = toset(["cluster", "node"])
  name_prefix        = "eks-${each.key}-role-${var.cluster_name}-"
  assume_role_policy = local.assume_role_policy[each.key]
}

resource "aws_iam_role_policy_attachment" "node" {
  for_each = toset([
    "AmazonEKSWorkerNodeMinimalPolicy",
    "AmazonEC2ContainerRegistryPullOnly",
  ])
  role       = aws_iam_role.aws_eks_cluster["node"].name
  policy_arn = "arn:aws:iam::aws:policy/${each.key}"
}

resource "aws_iam_role_policy_attachment" "cluster" {
  for_each = toset([
    "AmazonEKSClusterPolicy",
    "AmazonEKSComputePolicy",
    "AmazonEKSBlockStoragePolicyV2",
    "AmazonEKSLoadBalancingPolicy",
    "AmazonEKSNetworkingPolicy",
  ])
  role       = aws_iam_role.aws_eks_cluster["cluster"].name
  policy_arn = "arn:aws:iam::aws:policy/${each.key}"
}

resource "aws_eks_access_entry" "cluster_admin" {
  for_each      = toset(var.eks_cluster_admin_users)
  cluster_name  = aws_eks_cluster.this.name
  principal_arn = each.key
}

resource "aws_eks_access_policy_association" "cluster_admin" {
  for_each      = toset(var.eks_cluster_admin_users)
  cluster_name  = aws_eks_cluster.this.name
  policy_arn    = "arn:aws:eks::aws:cluster-access-policy/AmazonEKSClusterAdminPolicy"
  principal_arn = each.key

  access_scope {
    type = "cluster"
  }

  depends_on = [aws_eks_access_entry.cluster_admin]
}
