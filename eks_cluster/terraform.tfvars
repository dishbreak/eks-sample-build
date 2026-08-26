number_of_availability_zones = 3

# uncomment to set up test ec2 instances for the subnets
# provision_test_instances     = true

# useful if the cluster needs a new name between setup and teardown
# cluster_name = "eks-auto-1"

eks_cluster_admin_users         = ["arn:aws:iam::825573321580:user/vishal"]
provision_client_vpn_connection = true
