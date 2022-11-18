job "pack-yurt-info-discovery" {
  type = "batch"
  datacenters = [[ .my.datacenters  | toJson ]]
  [[ template "namespace" . ]]
  [[ template "region" . ]]

  periodic {
    cron = "@hourly"
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
        args = ["info", "discover"]
      }

      env {
        YURT_BACKEND = "redis"
        NOMAD_ADDR = "http://${attr.unique.network.ip-address}:4646"
      }

      template {
        data = <<EOT
{{$allocID := env "NOMAD_ALLOC_ID" -}}
REDIS_ADDR="{{range nomadService 1 $allocID "pack-yurt-info-redis"}}{{.Address}}:{{.Port}}{{end}}"
NOMAD_TOKEN="{{ with nomadVar "nomad/jobs/pack-yurt-info-discovery" }}{{ .token }}{{ end }}"
EOT
        destination = "local/env"
        env = true
      }
    }
  }
}
