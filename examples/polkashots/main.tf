provider "aws" {
  region = "eu-west-1"
}

module "validator" {
  source = "../../"

  instance_name     = var.instance_name
  ssh_key           = var.ssh_key
  additional_volume = var.additional_volume
  enable_polkashots = var.enable_polkashots
}
