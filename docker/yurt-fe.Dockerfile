FROM golang:alpine as stage0
RUN apk add yarn
RUN mkdir -p /go/yurt-tools
COPY ./ /go/yurt-tools
RUN cd /go/yurt-tools && \
        go mod vendor && \
        go get github.com/UnnoTed/fileb0x && \
        go generate ./... && \
        CGO_ENABLED=0 go build -o /yurt-fe cmd/yurt-fe/main.go

FROM alpine:latest
COPY --from=stage0 /yurt-fe /
CMD ["/yurt-fe"]
EXPOSE 8080/tcp
USER 1000
