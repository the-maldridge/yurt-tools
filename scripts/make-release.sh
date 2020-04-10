#!/bin/sh

if [ -z $VERSION ] ; then
    echo '$VERSION must be set!'
    exit 1
fi

for container in task-discover task-versions up2date yurt-fe trivy-dispatch ; do
    docker build -f docker/$container.Dockerfile -t yurttools/$container:$VERSION .
done
