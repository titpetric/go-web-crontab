FROM alpine:latest

ARG TARGETARCH

RUN apk update && apk add --no-cache bash coreutils curl
RUN which uname
RUN id nobody 2>/dev/null || adduser -D -H -u 65534 nobody
RUN mkdir -p /app /data && chown nobody:nobody /data && ln -s /data /app/data

COPY ./bin/webcron-linux-${TARGETARCH} /usr/local/bin/webcron
COPY ./cron.d /app/cron.d
COPY ./cron.scripts /app/cron.scripts

WORKDIR /app

EXPOSE 8080

USER nobody
ENTRYPOINT ["/usr/local/bin/webcron"]
