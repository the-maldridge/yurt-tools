job "trivy-server" {
  type = "service"
  datacenters = ["DC1"]
  group "app-scan" {
    task "trivy" {
      driver = "docker"
      config {
        image = "aquasec/trivy:0.4.4"
        args = ["server", "--listen", "0.0.0.0:1284"]
        port_map {
          vulndb = 1284
        }
      }
      resources {
        cpu = 1000
        memory = 512
        network {
          mbits = 10
          port "vulndb" {
            static = 1284
          }
        }
      }
      service {
        name = "trivy"
        port = "vulndb"
        check {
          name     = "alive"
          type     = "tcp"
          interval = "10s"
          timeout  = "2s"
        }
      }
    }
  }
}
