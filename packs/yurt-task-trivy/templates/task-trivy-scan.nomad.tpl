job "yurt-task-trivy-scan" {
  [[ template "region" . ]]
  datacenters = [[.my.datacenters|toStringList]]
  type = "batch"

  parameterized {
    payload       = "forbidden"
    meta_required = [
      "TRIVY_CONTAINER",
      "TRIVY_OUTPUT",
    ]
  }

  group "app" {
    count = 1

    network {
      mode = "bridge"
    }

    service {
      name = "pack-yurt-trivy-scan"
      port = 1

      connect {
        sidecar_service {
          proxy {
            upstreams {
              destination_name = "pack-yurt-trivy-server"
              local_bind_port = 4954
            }
          }
        }
      }
    }

    task "app" {
      driver = "docker"

      vault {
        policies = [ [[.my.vault_policy|quote]] ]
      }

      config {
        image = "ghcr.io/aquasecurity/trivy:[[ .my.trivy_version ]]"
        entrypoint = ["/bin/sh"]
        command = "/local/do-scan"
      }

      resources {
        memory = 1000
      }

      env {
        NOMAD_ADDR=[[.my.nomad_addr|quote]]
        CONSUL_HTTP_ADDR="http://${attr.unique.network.ip-address}:8500"
      }

      template {
        data = <<EOT
{{- with secret [[.my.vault_path_nomad|quote]] }}
NOMAD_TOKEN={{.Data.secret_id}}
{{- end }}
{{- with secret [[.my.vault_path_consul|quote]] }}
CONSUL_HTTP_TOKEN={{.Data.token}}
{{- end }}
EOT
        destination = "secrets/env"
        env = true
      }

      template {
        data = <<EOH
#!/bin/sh
apk add curl
trivy client --remote http://localhost:4954 \
        -f json -o trivy.json \
        {{ env "NOMAD_META_TRIVY_CONTAINER" }}

curl -X PUT \
        -H "Authorization: Bearer $CONSUL_HTTP_TOKEN" \
        -d @trivy.json \
        $CONSUL_HTTP_ADDR/v1/kv/{{ env "NOMAD_META_TRIVY_OUTPUT" }}
EOH
        destination = "local/do-scan"
        perms = "755"
        change_mode = "noop"
      }
    }
  }
}
