# yurt-task-discovery

This pack installs the basic task discovery components of the
yurt-tools.  It is typically consumed as a dependency of other packs
that act on the information that is discovered.

## Variables

- `job_name` (string) - The name to use as the job name which overrides using
  the pack name
- `datacenters` (list of strings) - A list of datacenters in the region which
  are eligible for task placement
- `region` (string) - The region where jobs will be deployed
