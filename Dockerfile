FROM golang:1.18-alpine3.16 AS builder

RUN apk add --update --no-cache ca-certificates make git curl

WORKDIR /usr/local/src/app

ARG GOPROXY

COPY go.mod go.sum ./
COPY api/go.* ./api/
RUN go mod download

COPY . .

RUN make build


FROM alpine:3.16.0

RUN apk add --update --no-cache ca-certificates tzdata bash curl

SHELL ["/bin/bash", "-c"]

COPY --from=builder /usr/local/src/app/build/* /usr/local/bin/

EXPOSE 8000 8001
CMD todobackend-go-kit --http-addr :${PORT:-8000} --public-url ${PUBLIC_URL}
