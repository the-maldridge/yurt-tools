# up2date

This tool provides a way to list all docker containers in your cluster
and check if they're up to date.  It does not currently support
non-docker registries.

Authentication to the registry can be provided by the environment
variables `UP2DATE_REGISTRY_USERNAME` and `UP2DATE_REGISTRY_PASSWORD`.
A `NOMAD_TOKEN` needs to be available with read-only credentials.
After configured, a web page will be served on port 8080 showing all
tasks running in the cluster and the relative versions of the
underlying containers.  Up to 5 versions newer than what is being run
will be shown in the output.
