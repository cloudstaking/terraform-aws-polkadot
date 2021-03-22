# terraform-aws-polkadot

Terraform module for provisioning ready-to-use _single node_ (or optionally _active-standby_) Kusama/Polkadot validators in AWS. Besides infrastructure (security group, instance, volume, etc), it also does:

- Pulls the latest snapshot from [Polkashots](https://polkashots.io)
- Creates a docker-compose with the [latest polkadot's release](https://github.com/paritytech/polkadot/releases) and nginx reverse-proxy (for libp2p port).

## Usage

```hcl
module "kusama_validator" {
  source = "github.com/cloudstaking/terraform-aws-polkadot?ref=1.0.0"

  instance_name       = "validator"
  ssh_key             = "ssh-rsa XXXXXXXXXXXXXX"
}
```

If `enable_polkashots` is set, it'll take ~10 minutes to download and extract the latest snapshot. You can check the process within the instance with `tail -f  /var/log/cloud-init-output.log`

<!-- BEGINNING OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
## Requirements

No requirements.

## Providers

| Name | Version |
|------|---------|
| aws | n/a |
| github | n/a |

## Modules

No Modules.

## Resources

| Name |
|------|
| [aws_ami](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/ami) |
| [aws_instance](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/instance) |
| [aws_key_pair](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/key_pair) |
| [aws_security_group](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/security_group) |
| [aws_vpc](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/vpc) |
| [github_release](https://registry.terraform.io/providers/integrations/github/latest/docs/data-sources/release) |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| chain | Chain name: kusama or polkadot. Variable required to download the latest snapshot from polkashots.io | `string` | `"kusama"` | no |
| disk\_size | Disk size. Because chain state constantly grows check the [requirements in the wiki](https://guide.kusama.network/docs/en/mirror-maintain-guides-how-to-validate-kusama) for the advisable sizes | `number` | `200` | no |
| enable\_polkashots | Pull latest Polkadot/Kusama (depending on chain variable) from polkashots.io | `bool` | `true` | no |
| instance\_name | Name of the instance | `string` | `"validator"` | no |
| instance\_type | Instance type: for Kusama m5.large is fine, for Polkadot maybe r5.2xlarge. This constantly change, check requirements section in the Polkadot wiki | `string` | `"m5.large"` | no |
| polkadot\_additional\_common\_flags | Application layer - when `enable_application_layer_docker = true`, the content of this variable will be appended to the polkadot command arguments | `string` | `""` | no |
| security\_group\_name | Security group name | `string` | `""` | no |
| security\_group\_whitelisted\_ssh\_ip | List of CIDRs the instance is going to accept SSH connections from. | `string` | `"0.0.0.0/0"` | no |
| ssh\_key | SSH Key to attach to the machine | `any` | n/a | yes |
| tags | A map of tags to add to all resources. | `map(string)` | `{}` | no |

## Outputs

| Name | Description |
|------|-------------|
| ssh | SSH command to connect to your validator |
| validator\_public\_ip | Validator public IP address, you can use it to SSH into it |
<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
