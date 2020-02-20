# up2date

This tool provides a way to list all docker containers in your cluster
and check if they're up to date.

Authentication to the registry can be provided by the environment
variables `<REGISTRY>_USERNAME` and `<REGISTRY>_PASSWORD`.  A
`NOMAD_TOKEN` needs to be available with read-only credentials.  After
configured, a web page will be served on port 8080 showing all tasks
running in the cluster and the relative versions of the underlying
containers.  Up to 5 versions newer than what is being run will be
shown in the output.

`<REGISTRY>` is specified as the registry name with underscores.
Currently recognized registries are:

  * `DOCKER_HUB`
  * `QUAY_IO`

Suggestions for new registries are welcome.
