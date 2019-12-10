FROM golang:latest as stage0
RUN mkdir -p /go/yurt-tools
COPY ./ /go/yurt-tools
RUN cd /go/yurt-tools && \
        go mod vendor && \
        CGO_ENABLED=0 go build -o /up2date cmd/up2date/main.go

FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
COPY --from=stage0 /up2date /
COPY --from=stage0 /go/yurt-tools/cmd/up2date/status.tpl /
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
CMD ["/up2date"]
EXPOSE 8080/tcp
USER 1000
