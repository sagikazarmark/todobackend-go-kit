FROM --platform=$BUILDPLATFORM tonistiigi/xx:1.6.1@sha256:923441d7c25f1e2eb5789f82d987693c47b8ed987c4ab3b075d6ed2b5d6779a3 AS xx

FROM --platform=$BUILDPLATFORM golang:1.23.5-alpine3.20@sha256:def59a601e724ddac5139d447e8e9f7d0aeec25db287a9ee1615134bcda266e2 AS builder

COPY --from=xx / /

RUN apk add --update --no-cache ca-certificates make git curl clang lld

ARG TARGETPLATFORM

RUN xx-apk --update --no-cache add musl-dev gcc

RUN xx-go --wrap

WORKDIR /usr/local/src/todobackend-go-kit

ARG GOPROXY

ENV CGO_ENABLED=1

COPY go.mod go.sum ./
COPY api/go.mod api/go.sum ./api/
RUN go mod download

COPY . .

RUN make build
RUN xx-verify build/todobackend-go-kit


FROM gcr.io/distroless/base-debian11:latest@sha256:ac69aa622ea5dcbca0803ca877d47d069f51bd4282d5c96977e0390d7d256455 AS distroless

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /usr/local/src/todobackend-go-kit/build/* /usr/local/bin/
COPY --from=builder /usr/local/src/todobackend-go-kit/go.* /usr/local/src/todobackend-go-kit/
COPY --from=builder /usr/local/src/todobackend-go-kit/api/go.* /usr/local/src/todobackend-go-kit/

EXPOSE 8000 8001
CMD todobackend-go-kit --http-addr :${PORT:-8000} --public-url ${PUBLIC_URL}

FROM redhat/ubi8-micro:8.10@sha256:22448ec2e9234d99a2cd9adf9e571a367d1e4b0e1546b2b5c36518e5183e1b32 AS ubi8

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /usr/local/src/todobackend-go-kit/build/* /usr/local/bin/
COPY --from=builder /usr/local/src/todobackend-go-kit/go.* /usr/local/src/todobackend-go-kit/
COPY --from=builder /usr/local/src/todobackend-go-kit/api/go.* /usr/local/src/todobackend-go-kit/

EXPOSE 8000 8001
CMD todobackend-go-kit --http-addr :${PORT:-8000} --public-url ${PUBLIC_URL}

FROM alpine:3.21.2@sha256:56fa17d2a7e7f168a043a2712e63aed1f8543aeafdcee47c58dcffe38ed51099 AS alpine

RUN apk add --update --no-cache ca-certificates tzdata bash

SHELL ["/bin/bash", "-c"]

COPY --from=builder /usr/local/src/todobackend-go-kit/build/* /usr/local/bin/
COPY --from=builder /usr/local/src/todobackend-go-kit/go.* /usr/local/src/todobackend-go-kit/
COPY --from=builder /usr/local/src/todobackend-go-kit/api/go.* /usr/local/src/todobackend-go-kit/

EXPOSE 8000 8001
CMD todobackend-go-kit --http-addr :${PORT:-8000} --public-url ${PUBLIC_URL}
