job "yurt-tools-versions" {
  datacenters = ["DC1"]
  type = "batch"
  periodic {
    cron = "0 0 * * * *"
    prohibit_overlap = true
  }
  group "version" {
    task "task-versions" {
      driver = "docker"
      vault {
        policies = ["yurt-tools"]
        change_mode = "noop"
      }
      config {
        image = "yurttools/task-versions:v0.3.0"
      }
      resources {
        cpu = 500
        memory = 64
      }
      template {
        data = <<EOH
CONSUL_HTTP_ADDR="http://${attr.unique.network.ip-address}:8500"
NOMAD_ADDR=http://nomad.service.consul:4646
{{- with secret "consul/creds/yurt-tools" }}
CONSUL_HTTP_TOKEN={{.Data.token}}
{{- end }}
{{- with secret "nomad/creds/yurt-tools" }}
NOMAD_TOKEN={{.Data.secret_id}}
{{- end }}
EOH
        destination = "secrets/env"
        perms = "400"
        env = true
      }
    }
  }
}
