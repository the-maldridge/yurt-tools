job "trivy-scan" {
  type = "batch"
  datacenters = ["DC1"]
  parameterized {
    payload       = "forbidden"
    meta_required = [
      "TRIVY_CONTAINER",
      "TRIVY_REGISTRY",
      "TRIVY_OUTPUT",
    ]
  }
  group "app-scan" {
    task "trivy" {
      driver = "docker"
      vault {
        policies = [
          "trivy",
          "yurt-tools",
        ]
        change_mode = "noop"
      }
      config {
        image = "aquasec/trivy:0.4.4"
        entrypoint = ["/local/do-scan"]
      }
      resources {
        cpu = 1000
        memory = 512
      }
      template {
        data = <<EOH
CONSUL_HTTP_ADDR=http://{{ env "attr.unique.network.ip-address"}}:8500
{{- with secret "consul/creds/yurt-tools" }}
CONSUL_HTTP_TOKEN={{.Data.token}}
{{- end }}
EOH
        destination = "secrets/consul_env"
        perms = "400"
        change_mode = "noop"
        env = true
      }
      template {
        data = <<EOH
#!/bin/sh
apk add curl
if [ -f /secrets/{{env "NOMAD_META_TRIVY_REGISTRY" }} ] ; then
        . /secrets/{{ env "NOMAD_META_TRIVY_REGISTRY" }}
fi
trivy client --remote http://trivy.service.consul:1284 \
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
      template {
        data = <<EOH
export TRIVY_AUTH_URL=https://registry.hub.docker.com
{{- with secret "secret/data/prod/trivy" }}
export TRIVY_USERNAME="{{.Data.data.docker_username}}"
export TRIVY_PASSWORD="{{.Data.data.docker_password}}"
{{- end }}
EOH
        destination = "secrets/docker-hub"
        perms = "400"
        change_mode = "noop"
      }
    }
  }
}
