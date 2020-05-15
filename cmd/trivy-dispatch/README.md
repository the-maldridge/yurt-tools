# trivy-dispatch

This task is meant to run on a schedule and launch parameterized tasks
for every docker job running in your cluster.  These tasks will scan
the docker image with Trivy and write the results back to Consul.

