FROM golang:latest as stage0
RUN mkdir -p /go/yurt-tools
COPY ./ /go/yurt-tools
RUN cd /go/yurt-tools && \
        go mod vendor && \
        CGO_ENABLED=0 go build -o /task-discover cmd/task-discover/main.go

FROM scratch
COPY --from=stage0 /task-discover /
CMD ["/task-discover"]
USER 1000
