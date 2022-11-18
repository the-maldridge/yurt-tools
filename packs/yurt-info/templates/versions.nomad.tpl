job "pack-yurt-info-versions" {
  type = "batch"
  datacenters = [[ .my.datacenters  | toJson ]]
  [[ template "namespace" . ]]
  [[ template "region" . ]]

  periodic {
    cron = "@daily"
  }
  
  group "app" {
    count = 1

    network {
      mode = "bridge"
    }

    task "discover" {
      driver = "docker"

      config {
        image = "ghcr.io/the-maldridge/yurt-tools:[[ .my.yurt_version ]]"
        args = ["info", "version-check"]
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
        env = true
      }
    }
  }
}
