output "repository_url" {
  value = "${coalesce("${module.ecr.repository_url}", "** unset **")}"
}
