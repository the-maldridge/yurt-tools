FROM golang:latest as stage0
RUN mkdir -p /go/yurt-tools
COPY ./ /go/yurt-tools
RUN cd /go/yurt-tools && \
        go mod vendor && \
        CGO_ENABLED=0 go build -o /trivy-dispatch cmd/trivy-dispatch/main.go

FROM scratch
COPY --from=stage0 /trivy-dispatch /
CMD ["/trivy-dispatch"]
USER 1000
