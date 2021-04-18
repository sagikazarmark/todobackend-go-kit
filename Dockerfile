FROM alpine:3.13.4 AS builder

RUN apk add --update --no-cache bash ca-certificates curl git build-base libc6-compat gcompat

# RUN cd /tmp; GOBIN=/build go get github.com/go-delve/delve/cmd/dlv

WORKDIR /usr/local/src/todobackend-go-kit

ARG PLZ_BUILD_CONFIG
ARG PLZ_OVERRIDES
ARG PLZ_CONFIG_PROFILE

ENV PLZ_ARGS="-p -o \"build.path:${PATH}\""

COPY .plzconfig* pleasew ./
RUN ./pleasew update

COPY BUILD .

COPY .build_defs ./.build_defs/
COPY tools ./tools/
RUN ./pleasew build //tools:go_toolchain

COPY third_party ./third_party/
RUN ./pleasew build //third_party/...

COPY . .

RUN ./pleasew export outputs -o /usr/local/bin //cmd/todo


FROM alpine:3.13.4

RUN apk add --update --no-cache ca-certificates tzdata bash curl libc6-compat

SHELL ["/bin/bash", "-c"]

# set up nsswitch.conf for Go's "netgo" implementation
# https://github.com/gliderlabs/docker-alpine/issues/367#issuecomment-424546457
RUN test ! -e /etc/nsswitch.conf && echo 'hosts: files dns' > /etc/nsswitch.conf

COPY --from=builder /usr/local/bin/* /usr/local/bin/

EXPOSE 8000 8001
CMD todo --http-addr :${PORT:-8000} --public-url ${PUBLIC_URL}
