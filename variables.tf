variable "security_group_name" {
  default     = ""
  description = "Security group name"
}

variable "security_group_whitelisted_ssh_ip" {
  default     = "0.0.0.0/0"
  description = "List of CIDRs the instance is going to accept SSH connections from."
}

variable "instance_name" {
  default     = "validator"
  description = "Name of the instance"
}

variable "instance_type" {
  default     = "m5.large"
  description = "Instance type: for Kusama m5.large is fine, for Polkadot maybe r5.2xlarge. This constantly change, check requirements section in the Polkadot wiki"
}

variable "disk_size" {
  description = "Disk size. Because chain state constantly grows check the [requirements in the wiki](https://guide.kusama.network/docs/en/mirror-maintain-guides-how-to-validate-kusama) for the advisable sizes"
  default     = 200
}

variable "chain" {
  description = "Chain name: kusama or polkadot. Variable required to download the latest snapshot from polkashots.io"
  default     = "kusama"
}

variable "enable_polkashots" {
  default     = true
  description = "Pull latest Polkadot/Kusama (depending on chain variable) from polkashots.io"
  type        = bool
}

variable "ssh_key" {
  description = "SSH Key to attach to the machine"
}

variable "tags" {
  description = "A map of tags to add to all resources."
  type        = map(string)
  default     = {}
}

variable "polkadot_additional_common_flags" {
  default     = ""
  description = "Application layer - when `enable_application_layer_docker = true`, the content of this variable will be appended to the polkadot command arguments"
}
