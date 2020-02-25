job "yurt-tools-trivy-dispatcher" {
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
        image = "yurttools/trivy-dispatch:v0.4.0"
      }
      resources {
        cpu = 500
        memory = 64
      }
      template {
        data = <<EOH
CONSUL_HTTP_ADDR=http://{{ env "attr.unique.network.ip-address" }}:8500
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
