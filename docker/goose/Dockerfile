FROM alpine:3.17.0

RUN apk update && \
    apk upgrade && \
    apk add bash && \
    rm -rf /var/cache/apk/*

ADD https://github.com/pressly/goose/releases/download/v3.7.0/goose_linux_x86_64 /bin/goose
ADD entrypoint.sh /migrations/

RUN chmod +x /bin/goose

ENTRYPOINT ["entrypoint.sh"]