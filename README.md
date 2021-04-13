# terraform-aws-polkadot

Terraform module to bootstrap ready-to-use _single node_ (or optionally _active-standby_) Kusama/Polkadot validators in AWS. Besides infrastructure components (security group, instance, volume, etc), it also:

- Optionally pulls latest snapshot from [Polkashots](https://polkashots.io)
- [Node exporter](https://github.com/prometheus/node_exporter) with HTTPs to securly pull metrics from your monitoring systems.
- Nginx as a reverse proxy for libp2p
- Support for different deplotments methods: either using docker/docker-compose or deploying the binary itself in the host.

It uses the latest official Ubuntu 20.04 LTS (no custom image). 

## Usage

```hcl
module "kusama_validator" {
  source = "github.com/cloudstaking/terraform-aws-polkadot?ref=1.1.0"

  instance_name = "validator"
  ssh_key       = "ssh-rsa XXXXXXXXXXXXXX"
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

## Modules

| Name | Source | Version |
|------|--------|---------|
| cloud_init | github.com/cloudstaking/terraform-cloudinit-polkadot?ref=main |  |

## Resources

| Name |
|------|
| [aws_ami](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/ami) |
| [aws_instance](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/instance) |
| [aws_key_pair](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/key_pair) |
| [aws_security_group](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/security_group) |
| [aws_vpc](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/vpc) |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| application\_layer | You can deploy the Polkadot using docker containers or in the host itself (using the binary) | `string` | `"host"` | no |
| chain | Chain name: kusama or polkadot. Variable required to download the latest snapshot from polkashots.io | `string` | `"kusama"` | no |
| disk\_size | Disk size. Because chain state constantly grows check the [requirements in the wiki](https://guide.kusama.network/docs/en/mirror-maintain-guides-how-to-validate-kusama) for the advisable sizes | `number` | `200` | no |
| enable\_polkashots | Pull latest Polkadot/Kusama (depending on chain variable) from polkashots.io | `bool` | `false` | no |
| http\_password | Password to access endpoints (e.g node\_exporter) | `string` | `""` | no |
| http\_username | Username to access endpoints (e.g node\_exporter) | `string` | `""` | no |
| instance\_name | Name of the instance | `string` | `"validator"` | no |
| instance\_type | Instance type: for Kusama m5.large is fine, for Polkadot maybe r5.2xlarge. This constantly change, check requirements section in the Polkadot wiki | `string` | `"m5.large"` | no |
| p2p\_port | P2P port for Polkadot service, used in `--listen-addr` args | `number` | `30333` | no |
| polkadot\_additional\_common\_flags | CLI arguments appended to the polkadot service (e.g validator name) | `string` | `""` | no |
| proxy\_port | nginx reverse-proxy port to expose Polkadot's libp2p port. Polkadot's libp2p port should not be exposed directly for security reasons (DOS) | `number` | `80` | no |
| public\_fqdn | Public domain for validator. If set, Caddy will use it to request LetsEncrypt certs. This variable is particulary useful to provide a secure channel (HTTPs) for [node\_exporter](https://github.com/prometheus/node_exporter) | `string` | `""` | no |
| security\_group\_name | Security group name | `string` | `""` | no |
| ssh\_key | SSH Key to attach to the machine | `any` | n/a | yes |
| tags | A map of tags to add to all resources. | `map(string)` | `{}` | no |

## Outputs

| Name | Description |
|------|-------------|
| http\_password | Password to access private endpoints (e.g node\_exporter) |
| http\_username | Username to access private endpoints (e.g node\_exporter) |
| ssh | SSH command to connect to your validator |
| validator\_public\_ip | Validator public IP address, you can use it to SSH into it |
<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
