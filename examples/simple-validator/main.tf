provider "aws" {
  region = "eu-west-1"
}

module "validator" {
  source = "../../"

  instance_name     = var.instance_name
  ssh_key           = var.ssh_key
  enable_polkashots = true

  polkadot_additional_common_flags = "--name=CLOUDSTAKING-TEST --telemetry-url 'wss://telemetry.polkadot.io/submit/ 1'"

  tags = {
    terraform = true
    leg       = "blue"
  }
}
