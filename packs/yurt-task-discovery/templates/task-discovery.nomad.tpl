job [[ template "job_name" . ]] {
  [[ template "region" . ]]
  datacenters = [[ .my.datacenters  | toJson ]]
  type = "batch"

  periodic {
    cron = [[.my.cronspec|quote]]
    prohibit_overlap = true
  }

  meta {
    # https://www.hashicorp.com/blog/running-duplicate-batch-jobs-in-hashicorp-nomad
    run_uuid = "${uuidv4()}"
  }

  group "app" {
    count = 1

    network {
      mode = "bridge"
    }

    task "discover" {
      driver = "docker"

      vault {
        policies = [ [[.my.vault_policy|quote]] ]
      }

      config {
        image = "ghcr.io/the-maldridge/yurt-tools:[[ .my.yurt_version ]]"
        args = ["info", "discover"]
      }

      env {
        NOMAD_ADDR=[[.my.nomad_addr|quote]]
        CONSUL_PREFIX=[[.my.consul_prefix|quote]]
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
    }
  }
}
