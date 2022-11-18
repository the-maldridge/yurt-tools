job "pack-yurt-info-trivy-scanner" {
  type = "batch"
  datacenters = [[ .my.datacenters  | toJson ]]
  [[ template "namespace" . ]]
  [[ template "region" . ]]

  parameterized {
    payload       = "forbidden"
    meta_required = [
      "TRIVY_CONTAINER",
      "TRIVY_OUTPUT",
    ]
    meta_optional = ["TRIVY_REGISTRY"]
  }
  
  group "app" {
    count = 1

    network {
      mode = "bridge"
    }

    task "scan" {
      driver = "docker"

      config {
        image = "ghcr.io/aquasecurity/trivy:[[ .my.trivy_version ]]"
        entrypoint = ["/bin/sh"]
        command = "/local/do-scan"
      }

      resources {
        memory = 1000
      }

      template {
        data = <<EOH
#!/bin/sh
trivy image --server http://{{ range nomadService 1 "a" "pack-yurt-trivy-server" }}{{.Address}}:{{.Port}}{{end}} \
        -f json -o /local/trivy.json \
        {{ env "NOMAD_META_TRIVY_CONTAINER" }}

apk add redis
redis-cli -u redis://{{ range nomadService 1 "a" "pack-yurt-info-redis" }}{{.Address}}:{{.Port}}{{end}} -x SET {{ env "NOMAD_META_TRIVY_OUTPUT" }} < /local/trivy.json 
EOH
        destination = "local/do-scan"
        perms = "755"
        change_mode = "noop"
      }
    }
  }
}
