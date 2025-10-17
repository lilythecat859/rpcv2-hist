# AGPL-3.0
# multi-stage distroless build – no shell, no cgo, static binary
FROM golang:1.22-alpine AS builder
RUN apk add --no-cache git make
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 make build

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /src/bin/rpcv2-hist /rpcv2-hist
EXPOSE 8899 8080 9090 9091
ENTRYPOINT ["/rpcv2-hist"]