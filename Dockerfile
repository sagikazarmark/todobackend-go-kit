FROM alpine:3.14.3 AS builder

RUN apk add --update --no-cache bash ca-certificates curl git build-base libc6-compat

# RUN cd /tmp; GOBIN=/build go get github.com/go-delve/delve/cmd/dlv

WORKDIR /usr/local/src/todobackend-go-kit

ARG PLZ_BUILD_CONFIG
ARG PLZ_OVERRIDES
ARG PLZ_CONFIG_PROFILE

ENV PLZ_ARGS="-p -o \"build.path:${PATH}\""

COPY .plzconfig* pleasew ./
RUN ./pleasew update

COPY BUILD .

# Required for certain please functions to work
RUN git init && git config --global user.email "you@example.com" && git config --global user.name "Your Name" && git commit -m 'dummy' --allow-empty

COPY tools ./tools/
RUN ./pleasew build //tools:go_toolchain

COPY third_party ./third_party/
RUN ./pleasew build //third_party/...

RUN rm -rf .git

COPY . .

RUN ./pleasew export outputs -o /usr/local/bin :todobackend


FROM alpine:3.14.3

RUN apk add --update --no-cache ca-certificates tzdata bash curl libc6-compat

SHELL ["/bin/bash", "-c"]

# set up nsswitch.conf for Go's "netgo" implementation
# https://github.com/gliderlabs/docker-alpine/issues/367#issuecomment-424546457
RUN test ! -e /etc/nsswitch.conf && echo 'hosts: files dns' > /etc/nsswitch.conf

COPY --from=builder /usr/local/bin/* /usr/local/bin/

EXPOSE 8000 8001
CMD todobackend --http-addr :${PORT:-8000} --public-url ${PUBLIC_URL}
