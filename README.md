# terraform-scaleway-polkadot

Terraform module for provisioning ready-to-use _single node_ (or optionally _active-standby_) Kusama/Polkadot validators in AWS. Besides infrastructure (security group, instance, volume, etc), it also does:

- Pulls the latest snapshot from [Polkashots](https://polkashots.io)
- Creates a docker-compose with the [latest polkadot's release](https://github.com/paritytech/polkadot/releases) and nginx reverse-proxy (for libp2p port).

## Usage

```hcl
module "kusama_validator" {
  source = "github.com/cloudstaking/terraform-scaleway-polkadot?ref=1.2.0"

  instance_name       = "ksm-validator"
  security_group_name = "ksm-validator"
  chain               = "kusama"
}
```

If `enable_polkashots` is set, it'll take ~10 minutes to download and extract the latest snapshot. You can check the process within the instance

```sh
$ tail -f  /var/log/cloud-init-output.log
4323850K .......... .......... .......... .......... .......... 99%  313M 0s
4323900K .......... ......                                     100% 56.2M=78s

2021-03-17 14:27:38 (54.2 MB/s) - ‘/srv/kusama.RocksDb.7z’ saved [4427690609/4427690609]

Scanning the drive for archives:
1 file, 4427690609 bytes (4223 MiB)

Extracting archive: kusama.RocksDb.7z
```

<!-- BEGINNING OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
## Requirements

No requirements.

## Providers

| Name | Version |
|------|---------|
| github | n/a |
| scaleway | n/a |

## Modules

No Modules.

## Resources

| Name |
|------|
| [github_release](https://registry.terraform.io/providers/integrations/github/latest/docs/data-sources/release) |
| [scaleway_instance_ip](https://registry.terraform.io/providers/scaleway/scaleway/latest/docs/resources/instance_ip) |
| [scaleway_instance_security_group](https://registry.terraform.io/providers/scaleway/scaleway/latest/docs/resources/instance_security_group) |
| [scaleway_instance_server](https://registry.terraform.io/providers/scaleway/scaleway/latest/docs/resources/instance_server) |
| [scaleway_instance_volume](https://registry.terraform.io/providers/scaleway/scaleway/latest/docs/resources/instance_volume) |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| additional\_volume | By default, DEV1-M comes with 40GB disk size. Set this variable in order to create an additional volume (mounted in /srv) | `bool` | `true` | no |
| additional\_volume\_size | Volume size where the chain state is going to be saved (only applies if additional\_volume variable set) - check Kusama/Polkadot requirements | `number` | `200` | no |
| chain | Chain name: kusama or polkadot. Variable required to download the latest snapshot from polkashots.io | `string` | `"kusama"` | no |
| enable\_docker\_compose | Application layer - create a docker-compose.yml (`/srv/docker-compose.yml`) with the latest polkadot version and nginx as a reverse-proxy | `bool` | `false` | no |
| enable\_polkashots | Pull latest Polkadot/Kusama (depending on chain variable) from polkashots.io | `bool` | `true` | no |
| instance\_name | Name of the Scaleway instance | `string` | `"validator"` | no |
| instance\_type | Instance type: for Kusama DEV1-M is fine, for Polkadot maybe GP1-M. Check requirements in the Polkadot wiki | `string` | `"DEV1-M"` | no |
| polkadot\_additional\_common\_flags | Application layer - when `enable_docker_compose = true`, the content of this variable will be appended to the polkadot command arguments | `string` | `""` | no |
| security\_group\_name | Security group name | `string` | `"validator-sg"` | no |
| security\_group\_whitelisted\_ssh\_ip | List of CIDRs the instance is going to accept SSH connections from. | `string` | `"0.0.0.0/0"` | no |
| tags | Tags for the instance | `list` | `[]` | no |

## Outputs

| Name | Description |
|------|-------------|
| validator\_public\_ip | Validator public IP address, you can use it to SSH into it |
<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
