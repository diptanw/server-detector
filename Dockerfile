FROM golang:1.14-alpine3.11 AS builder

ARG package

WORKDIR /src/
COPY . .

RUN apk add --update --no-cache git
RUN go build -o srv -v -ldflags "-s -w" ./cmd/${package}

FROM alpine:3.11

RUN apk upgrade --update --no-cache \
	&& addgroup -S 65011 \
	&& adduser -D -S -G 65011 65011

USER 65011:65011

COPY --from=builder /src/srv /usr/local/bin/
COPY --from=builder /src/api/openapi.yml /usr/local/bin/

EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=3s --retries=5 \
	CMD wget -q -O - http://localhost:8080/health || exit 1

ENTRYPOINT [ "srv" ]

