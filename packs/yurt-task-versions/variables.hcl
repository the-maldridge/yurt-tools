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

variable "datacenters" {
  description = "A list of datacenters in the region which are eligible for task placement"
  type        = list(string)
}

variable "cronspec" {
  description = "A valid cron expression that defines when new tasks should be picked up from nomad."
  type = string
  default = "@daily"
}

variable "consul_prefix" {
  description = "The prefix that all task components will share for access to consul"
  type = string
  default = "yurt-tools/tasks"
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
