FROM golang:latest as stage0
RUN mkdir -p /go/yurt-tools
COPY ./ /go/yurt-tools
RUN cd /go/yurt-tools && \
        go mod vendor && \
        CGO_ENABLED=0 go build -o /version-checker cmd/version-checker/main.go

FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM alpine:latest
COPY --from=stage0 /version-checker /
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
CMD ["/version-checker"]
USER 1000
