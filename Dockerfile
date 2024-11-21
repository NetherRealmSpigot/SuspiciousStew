ARG CODE="Stew"
FROM golang:1.24-alpine3.21 AS build-env
ARG CODE

# Build deps
RUN apk --no-cache add build-base git && rm -rf /var/cache/apk/*

WORKDIR /build

COPY . .
RUN CODE=$CODE make build


FROM alpine:3.21
ARG CODE

RUN <<EOS
(deluser --remove-home xfs 2>/dev/null || true)
(deluser --remove-home www-data 2>/dev/null || true)
(delgroup www-data 2>/dev/null || true)
(delgroup xfs 2>/dev/null || true)
addgroup -S -g 33 www-data
adduser -S -D -u 33 -s /sbin/nologin -h /var/www -G www-data www-data
EOS

WORKDIR /app

COPY --from=build-env --chown=www-data:www-data --chmod=755 /build/$CODE . 
RUN if [ "$CODE" != "Stew" ]; then ln -s $CODE Stew; fi
COPY --chown=www-data:www-data --chmod=755 docker-entrypoint.sh /
RUN apk --no-cache add tini && rm -rf /var/cache/apk/*

ENV STEWAPI_LISTEN_ADDRESS=0.0.0.0
ENV STEWAPI_LISTEN_PORT=8080

USER www-data
EXPOSE 8080

ENTRYPOINT ["tini", "--", "/docker-entrypoint.sh"]