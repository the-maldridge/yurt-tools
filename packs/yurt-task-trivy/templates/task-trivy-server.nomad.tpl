job "yurt-task-trivy-server" {
  [[ template "region" . ]]
  datacenters = [[.my.datacenters|toStringList]]
  type = "service"

  group "server" {
    count = 1

    network {
      mode = "bridge"
    }

    service {
      name = "pack-yurt-trivy-server"
      port = 4954

      connect {
        sidecar_service {}
      }
    }
    
    task "app" {
      driver = "docker"

      config {
        image = "ghcr.io/aquasecurity/trivy:[[ .my.trivy_version ]]"
        command = "server"
      }

      resources {
        memory = 600
      }
    }
  }
}
