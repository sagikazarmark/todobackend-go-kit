# Build image
FROM golang:1.15-alpine AS builder

ENV GOFLAGS="-mod=readonly"
ENV CGO_ENABLED=0

RUN apk add --update --no-cache bash ca-certificates curl git mercurial

RUN cd /tmp; GOBIN=/build go get github.com/go-delve/delve/cmd/dlv

RUN mkdir -p /workspace
WORKDIR /workspace

ARG GOPROXY

COPY go.mod go.sum ./
COPY api/go.mod api/go.sum ./api/
RUN go mod download

COPY . .

ARG PLZ_BUILD_CONFIG
ARG PLZ_OVERRIDES
ARG PLZ_CONFIG_PROFILE

RUN echo -e "[build]\npath = ${PATH}" > .plzconfig.local

RUN ./pleasew -p export outputs -o /build //cmd/...


# Final image
FROM alpine:3.12

RUN apk add --update --no-cache ca-certificates tzdata bash curl libc6-compat

SHELL ["/bin/bash", "-c"]

# set up nsswitch.conf for Go's "netgo" implementation
# https://github.com/gliderlabs/docker-alpine/issues/367#issuecomment-424546457
RUN test ! -e /etc/nsswitch.conf && echo 'hosts: files dns' > /etc/nsswitch.conf

COPY --from=builder /build/* /usr/local/bin/

EXPOSE 8000 8001
CMD todo --http-addr :${PORT:-8000} --public-url ${PUBLIC_URL}
