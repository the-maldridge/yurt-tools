job "pack-yurt-info" {
  type = "service"
  datacenters = [[ .my.datacenters  | toJson ]]
  [[ template "namespace" . ]]
  [[ template "region" . ]]

  group "redis" {
    count = 1

    network {
      mode = "bridge"
      port "redis" { to = 6379 }
    }

    service {
      name = "pack-yurt-info-redis"
      port = "redis"
      provider = "nomad"
    }

    task "redis" {
      driver = "docker"

      config {
        image = "docker.io/library/redis:7"
      }
    }
  }

  group "frontend" {
    count = 1

    network {
      mode = "bridge"
      port "http" { to = 8080 }
    }

    service {
      name = "pack-yurt-info"
      port = "http"
      provider = [[ .my.service_provider | quote ]]
      tags = [[ .my.service_tags | toJson ]]
    }

    task "frontend" {
      driver = "docker"

      config {
        image = "ghcr.io/the-maldridge/yurt-tools:[[ .my.yurt_version ]]"
        args = ["http", "taskinfo"]
      }

      env {
        YURT_BACKEND = "redis"
      }

      template {
        data = <<EOT
{{$allocID := env "NOMAD_ALLOC_ID" -}}
REDIS_ADDR="{{range nomadService 1 $allocID "pack-yurt-info-redis"}}{{.Address}}:{{.Port}}{{end}}"
EOT
        destination = "local/env"
        change_mode = "restart"
        env = true
      }
    }
  }

  group "trivy-server" {
    count = 1

    network {
      mode = "bridge"
      port "trivy" { to = 4965 }
    }

    service {
      name = "pack-yurt-trivy-server"
      provider = "nomad"
      port = "trivy"
    }
    
    task "app" {
      driver = "docker"

      config {
        image = "ghcr.io/aquasecurity/trivy:[[ .my.trivy_version ]]"
        command = "server"
        args = ["--listen", "0.0.0.0:4965"]
      }

      resources {
        memory = 600
      }
    }
  }
}
