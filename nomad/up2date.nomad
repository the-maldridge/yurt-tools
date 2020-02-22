job "up2date" {
  datacenters = ["DC1"]

  group "up2date" {
    task "up2date" {
      driver = "docker"
      vault {
        policies = ["up2date"]
        change_mode = "noop"
      }
      config {
        image = "yurttools/up2date:v0.2.0"
        port_map {
          http = 8080
        }
      }
      resources {
        cpu    = 500
        memory = 128
        network {
          mbits = 10
          port "http" {}
        }
      }
      template {
        data = <<EOH
NOMAD_ADDR=http://nomad.service.consul:4646
{{- with secret "nomad/creds/up2date" }}
NOMAD_TOKEN="{{.Data.secret_id}}"
{{- end }}
{{- with secret "secret/data/prod/up2date" }}
UP2DATE_REGISTRY_USERNAME="{{.Data.data.username}}"
UP2DATE_REGISTRY_PASSWORD="{{.Data.data.password}}"
{{- end }}
EOH
        destination = "secrets/env"
        perms = "644"
        env = true
      }

      service {
        name = "up2date"
        port = "http"
        check {
          name     = "alive"
          type     = "http"
          path     = "/"
          interval = "10s"
          timeout  = "2s"
        }
      }
    }
  }
}
