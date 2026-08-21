provider "aws" {
  default_tags {
    tags = {
      "project" = "eks-sample-build"
    }
  }
  region = "us-west-2"
}
