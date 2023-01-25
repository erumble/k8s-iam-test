locals {
  some_local = "some_value"
}

output "some_output" {
  value = local.some_local
}
