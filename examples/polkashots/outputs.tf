output "public_ip" {
  value = module.validator.validator_public_ip
}

output "ssh_cmd" {
  value = module.validator.ssh
}
