FROM --platform=$BUILDPLATFORM tonistiigi/xx:1.6.1@sha256:923441d7c25f1e2eb5789f82d987693c47b8ed987c4ab3b075d6ed2b5d6779a3 AS xx

FROM --platform=$BUILDPLATFORM golang:1.23.4-alpine3.20@sha256:9a31ef0803e6afdf564edc8ba4b4e17caed22a0b1ecd2c55e3c8fdd8d8f68f98 AS builder

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

FROM redhat/ubi8-micro:8.8@sha256:c743e8d6f673f8287a07e3590cbf65dfa7c5c21bb81df6dbd4d9a2fcf21173cd AS ubi8

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /usr/local/src/todobackend-go-kit/build/* /usr/local/bin/
COPY --from=builder /usr/local/src/todobackend-go-kit/go.* /usr/local/src/todobackend-go-kit/
COPY --from=builder /usr/local/src/todobackend-go-kit/api/go.* /usr/local/src/todobackend-go-kit/

EXPOSE 8000 8001
CMD todobackend-go-kit --http-addr :${PORT:-8000} --public-url ${PUBLIC_URL}

FROM alpine:3.21.0@sha256:e323a465c03a31ad04374fc7239144d0fd4e2b92da6e3e0655580476d3a84621 AS alpine

RUN apk add --update --no-cache ca-certificates tzdata bash

SHELL ["/bin/bash", "-c"]

COPY --from=builder /usr/local/src/todobackend-go-kit/build/* /usr/local/bin/
COPY --from=builder /usr/local/src/todobackend-go-kit/go.* /usr/local/src/todobackend-go-kit/
COPY --from=builder /usr/local/src/todobackend-go-kit/api/go.* /usr/local/src/todobackend-go-kit/

EXPOSE 8000 8001
CMD todobackend-go-kit --http-addr :${PORT:-8000} --public-url ${PUBLIC_URL}
