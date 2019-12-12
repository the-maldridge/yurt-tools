# yurt-tools

Yurt Tools is a collection of tools for your Nomad cluster.  This
contains things that shouldn't necessarily be built in to Nomad, but
are very useful to have.

## [up2date](cmd/up2date/)

This tool is an update checker that can tell you if your docker
containers are out of date.  Useful if you pull unmodified containers
and don't have another mechanism for tracking upstreams.
