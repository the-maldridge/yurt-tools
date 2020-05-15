# task-discover

This task should run on a schedule and will scan your Nomad cluster
for jobs.  This task must run with permissions sufficient to List jobs
and sufficient to Read jobs.

The Nomad token that is used to contact the cluster will be read from
any location supported by the Nomad SDK.  It is recommended to use an
environment variable.
