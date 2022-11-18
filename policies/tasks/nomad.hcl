namespace "*" {
  capabilities = [
    "list-jobs",
    "read-job",
  ]
}

namespace "default" {
  capabilities = [
    "dispatch-job",
    "list-jobs",
    "read-job",
  ]
}
