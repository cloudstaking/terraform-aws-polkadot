locals {
  chain = {
    kusama   = { name = "kusama", short = "ksm" },
    polkadot = { name = "polkadot", short = "dot" }
    other    = { name = var.chain, short = var.chain }
  }

  security_group_name = var.security_group_name != "" ? var.security_group_name : "${var.instance_name}-sg"

  docker_compose = templatefile("${path.module}/templates/generate-docker-compose.sh.tpl", {
    chain                   = var.chain
    enable_polkashots       = var.enable_polkashots
    latest_version          = data.github_release.polkadot.release_tag
    additional_common_flags = var.polkadot_additional_common_flags
  })

  cloud_init = templatefile("${path.module}/templates/cloud-init.yaml.tpl", {
    chain             = lookup(local.chain, var.chain, local.chain.other)
    enable_polkashots = var.enable_polkashots
    additional_volume = var.additional_volume
    docker_compose    = base64encode(local.docker_compose)
  })
}

data "aws_vpc" "selected" {
  default = true
}

resource "aws_security_group" "validator" {
  name        = local.security_group_name
  description = "${var.instance_name} default security group"
  vpc_id      = data.aws_vpc.selected.id

  ingress {
    description = "SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = [var.security_group_whitelisted_ssh_ip]
  }
   
  ingress {
    description = "nginx (reverse-proxy for p2p port)"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
   
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = var.tags
}

data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Canonical
}

resource "aws_key_pair" "validator" {
  key_name   = var.instance_name
  public_key = var.ssh_key
}

resource "aws_instance" "validator" {
  ami                         = data.aws_ami.ubuntu.id
  instance_type               = var.instance_type
  associate_public_ip_address = true
  vpc_security_group_ids      = [aws_security_group.validator.id]
  key_name                    = aws_key_pair.validator.key_name
  user_data                   = local.cloud_init

  dynamic "ebs_block_device" {
    for_each = range(var.additional_volume ? 1 : 0)

    content {
      device_name = "/dev/sdb"
      volume_size = var.additional_volume_size
      volume_type = "gp2"
      delete_on_termination = true
    }
  }

  tags = merge(
    var.tags,
    {
      Name = var.instance_name
    },
  )
}

data "github_release" "polkadot" {
  repository  = "polkadot"
  owner       = "paritytech"
  retrieve_by = "latest"
}
