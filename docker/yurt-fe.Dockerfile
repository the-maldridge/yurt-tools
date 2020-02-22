FROM golang:latest as stage0
RUN mkdir -p /go/yurt-tools
COPY ./ /go/yurt-tools
RUN cd /go/yurt-tools && \
        go mod vendor && \
        CGO_ENABLED=0 go build -o /yurt-fe cmd/yurt-fe/main.go

FROM alpine:latest
COPY --from=stage0 /yurt-fe /
COPY --from=stage0 /go/yurt-tools/cmd/yurt-fe/static /static
COPY --from=stage0 /go/yurt-tools/cmd/yurt-fe/tmpl /tmpl
CMD ["/yurt-fe"]
EXPOSE 8080/tcp
USER 1000
