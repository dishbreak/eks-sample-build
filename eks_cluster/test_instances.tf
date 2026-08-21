data "aws_ami" "al23_latest" {
  most_recent = true
  owners = ["amazon"]
  
  filter {
    name = "name"
    values = ["al2023-ami-2023.*-x86_64"]
  }

  filter {
    name = "state"
    values = ["available"]
  }
  
  filter {
    name = "architecture"
    values = ["x86_64"]
  }

  filter {
    name = "root-device-type"
    values = ["ebs"]
  }

  filter {
    name = "virtualization-type"
    values = ["hvm"]
  }
}

resource "aws_instance" "test_inst" {
  for_each = aws_subnet.this
  instance_type = "t3.micro"
  subnet_id = each.value.id
  root_block_device {
    delete_on_termination = true
    volume_type = "gp3"
    volume_size = "10"
  }
  ami = data.aws_ami.al23_latest.id
  iam_instance_profile = aws_iam_instance_profile.inst.name

  # avoid churn on latest ami
  lifecycle {
    ignore_changes = [ ami ]
  }

  vpc_security_group_ids = [aws_security_group.inst.id]

  tags = {
    Name = "${each.key}-test-inst"
  }
}

data "aws_iam_policy_document" "inst_assume_role_policy" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "inst" {
  assume_role_policy = data.aws_iam_policy_document.inst_assume_role_policy.json
  name_prefix = "test-inst-role-"
}

resource "aws_iam_role_policy_attachment" "inst_ssm" {
  role = aws_iam_role.inst.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

resource "aws_iam_instance_profile" "inst" {
  role = aws_iam_role.inst.name
}

resource "aws_security_group" "inst" {
  name_prefix = "test_inst"
  vpc_id = aws_vpc.this.id
}

resource "aws_security_group_rule" "inst_outbound" {
  security_group_id = aws_security_group.inst.id
  type = "egress"
  protocol = "-1"
  from_port = 0
  to_port = 65536
  cidr_blocks = [ "0.0.0.0/0" ]
}
