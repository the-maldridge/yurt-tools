job "yurt-tools-frontend" {
  datacenters = ["DC1"]
  type = "service"
  group "frontend" {
    task "yurt-fe" {
      driver = "docker"
      vault {
        policies = ["yurt-tools"]
        change_mode = "noop"
      }
      config {
        image = "yurttools/yurt-fe:v0.3.0"
        port_map {
          http = 8080
        }
      }
      resources {
        cpu = 500
        memory = 64
        network {
          mbits = 10
          port "http" {}
        }
      }
      template {
        data = <<EOH
CONSUL_HTTP_ADDR="http://${attr.unique.network.ip-address}:8500"
{{- with secret "consul/creds/yurt-tools" }}
CONSUL_HTTP_TOKEN={{.Data.token}}
{{- end }}
EOH
        destination = "secrets/env"
        perms = "400"
        env = true
      }
      service {
        name = "yurt-tools"
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

