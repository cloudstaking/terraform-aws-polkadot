output "validator_public_ip" {
  value       = aws_instance.validator.public_ip
  description = "Validator public IP address, you can use it to SSH into it"
}

output "ssh" {
  value       = "ssh ubuntu@${aws_instance.validator.public_ip}"
  description = "SSH command to connect to your validator"
}
