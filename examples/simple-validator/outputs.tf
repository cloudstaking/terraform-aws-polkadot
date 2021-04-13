output "public_ip" {
  value = module.validator.validator_public_ip
}

output "ssh_cmd" {
  value = module.validator.ssh
}

output "http_username" {
  value = module.validator.http_username
}

output "http_password" {
  value = module.validator.http_password
}
