FROM --platform=$BUILDPLATFORM tonistiigi/xx:1.3.0@sha256:904fe94f236d36d65aeb5a2462f88f2c537b8360475f6342e7599194f291fb7e AS xx

FROM --platform=$BUILDPLATFORM golang:1.20.4-alpine3.16@sha256:6469405d7297f82d56195c90a3270b0806ef4bd897aa0628477d9959ab97a577 AS builder

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


FROM gcr.io/distroless/base-debian11:latest@sha256:73deaaf6a207c1a33850257ba74e0f196bc418636cada9943a03d7abea980d6d AS distroless

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

FROM alpine:3.18.3@sha256:7144f7bab3d4c2648d7e59409f15ec52a18006a128c733fcff20d3a4a54ba44a AS alpine

RUN apk add --update --no-cache ca-certificates tzdata bash

SHELL ["/bin/bash", "-c"]

COPY --from=builder /usr/local/src/todobackend-go-kit/build/* /usr/local/bin/
COPY --from=builder /usr/local/src/todobackend-go-kit/go.* /usr/local/src/todobackend-go-kit/
COPY --from=builder /usr/local/src/todobackend-go-kit/api/go.* /usr/local/src/todobackend-go-kit/

EXPOSE 8000 8001
CMD todobackend-go-kit --http-addr :${PORT:-8000} --public-url ${PUBLIC_URL}
