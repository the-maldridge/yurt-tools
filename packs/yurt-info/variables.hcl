variable "job_name" {
  description = "The name to use as the job name which overrides using the pack name"
  type        = string
  // If "", the pack name will be used
  default = ""
}

variable "region" {
  description = "The region where jobs will be deployed"
  type        = string
  default     = ""
}

variable "namespace" {
  description = "Namespace to deploy into"
  type = string
  default = ""
}

variable "datacenters" {
  description = "Datacenters to deploy to"
  type = list(string)
  default = ["dc1"]
}

variable "service_tags" {
  description = "List of tags to apply to the registered service"
  type = list(string)
  default = ["traefik.enable=true", "traefik.consulcatalog.connect=true"]
}

variable "service_provider" {
  description = "Provider of the service discovery layer (nomad or consul)"
  type = string
  default = "nomad"
}

variable "vault_policy" {
  description = "The name of the policy that grants access to a nomad token and consul token"
  type = string
  default = "pack-yurt-tasks"
}

variable "vault_path_consul" {
  description = "Path to the authorized consul secrets engine role in Vault"
  type = string
  default = "consul/creds/pack-yurt-tasks"
}

variable "yurt_version" {
  description = "Version of the yurt-tools to deploy"
  type = string
  default = "v0.5.0"
}

variable "trivy_version" {
  description = "Version of trivy to deploy"
  type = string
  default = "0.34.0"
}
