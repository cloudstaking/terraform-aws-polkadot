locals {
  security_group_name = var.security_group_name != "" ? var.security_group_name : "${var.instance_name}-sg"
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
    cidr_blocks = ["0.0.0.0/0"]
  }

  # https for TLS-ALPN challange only when public_fqdn is given:
  # https://caddyserver.com/docs/automatic-https#tls-alpn-challenge
  dynamic "ingress" {
    for_each = range(var.public_fqdn != "" ? 1 : 0)
    content {
      description = "https for TLS-ALPN challange when public_fqdn given"
      from_port   = 443
      to_port     = 443
      protocol    = "tcp"
      cidr_blocks = ["0.0.0.0/0"]
    }
  }

  ingress {
    description = "nginx (reverse-proxy for p2p port)"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "node_exporter"
    from_port   = 9100
    to_port     = 9100
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "polkadot exporter"
    from_port   = 9616
    to_port     = 9616
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
  user_data                   = module.cloud_init.clout_init

  root_block_device {
    volume_size = var.disk_size
  }

  tags = merge(
    var.tags,
    {
      Name = var.instance_name
    },
  )
}

module "cloud_init" {
  source = "github.com/cloudstaking/terraform-cloudinit-polkadot?ref=main"

  application_layer                = var.application_layer
  additional_volume                = false
  cloud_provider                   = "aws"
  chain                            = var.chain
  polkadot_additional_common_flags = var.polkadot_additional_common_flags
  enable_polkashots                = var.enable_polkashots
  p2p_port                         = var.p2p_port
  proxy_port                       = var.proxy_port
  public_fqdn                      = var.public_fqdn
  http_username                    = var.http_username
  http_password                    = var.http_password
}
