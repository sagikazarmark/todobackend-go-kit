FROM --platform=$BUILDPLATFORM tonistiigi/xx:1.2.1 AS xx

FROM --platform=$BUILDPLATFORM golang:1.20.3-alpine3.16@sha256:29c4e6e307eac79e5db29a261b243f27ffe0563fa1767e8d9a6407657c9a5f08 AS builder

COPY --from=xx / /

RUN apk add --update --no-cache ca-certificates make git curl clang lld

ARG TARGETPLATFORM

RUN xx-apk --update --no-cache add musl-dev gcc

RUN xx-go --wrap

WORKDIR /usr/local/src/todobackend-go-kit

ARG GOPROXY

ENV CGO_ENABLED=1

COPY go.mod go.sum ./
COPY api/go.* ./api/
RUN go mod download

COPY . .

RUN make build
RUN xx-verify build/todobackend-go-kit


FROM alpine:3.17.3@sha256:124c7d2707904eea7431fffe91522a01e5a861a624ee31d03372cc1d138a3126

RUN apk add --update --no-cache ca-certificates tzdata bash curl

SHELL ["/bin/bash", "-c"]

COPY --from=builder /usr/local/src/todobackend-go-kit/build/* /usr/local/bin/
COPY --from=builder /usr/local/src/todobackend-go-kit/go.* /usr/local/src/todobackend-go-kit/
COPY --from=builder /usr/local/src/todobackend-go-kit/api/go.* /usr/local/src/todobackend-go-kit/

EXPOSE 8000 8001
CMD todobackend-go-kit --http-addr :${PORT:-8000} --public-url ${PUBLIC_URL}
