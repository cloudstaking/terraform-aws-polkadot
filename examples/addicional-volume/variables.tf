variable "additional_volume" {
  type = bool
}

variable "instance_name" {
  description = "Name of the Scaleway instance"
}

variable "ssh_key" {
  description = "SSH Key to attach to the machine"
}

variable "tags" {
  description = "A map of tags to add to all resources."
  type        = map(string)
  default     = {}
}
