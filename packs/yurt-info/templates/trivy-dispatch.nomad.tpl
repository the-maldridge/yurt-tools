job "pack-yurt-info-trivy-dispatcher" {
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
        args = ["info", "trivy-dispatch"]
      }

      env {
        YURT_BACKEND = "redis"
        YURT_TRIVY_DISPATCHABLE = "pack-yurt-info-trivy-scanner"
        NOMAD_ADDR = "http://${attr.unique.network.ip-address}:4646"
      }

      template {
        data = <<EOT
{{$allocID := env "NOMAD_ALLOC_ID" -}}
REDIS_ADDR="{{range nomadService 1 $allocID "pack-yurt-info-redis"}}{{.Address}}:{{.Port}}{{end}}"
NOMAD_TOKEN="{{ with nomadVar "nomad/jobs/pack-yurt-info-trivy-dispatcher" }}{{ .token }}{{ end }}"
EOT
        destination = "local/env"
        env = true
      }
    }
  }
}
