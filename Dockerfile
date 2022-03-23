FROM golang:1.17-alpine as build
WORKDIR /yurt-tools
COPY . .
RUN go mod vendor && \
        CGO_ENABLED=0 go build -o /yurt ./cmd/yurt/main.go

FROM alpine:latest as cacerts
RUN apk add --no-cache ca-certificates

FROM scratch
COPY --from=build /yurt /yurt
COPY --from=cacerts /etc/ssl /etc/ssl
COPY theme /theme
ENTRYPOINT ["/yurt"]
