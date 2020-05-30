FROM golang:alpine

RUN apk update --no-cache && \
    apk upgrade --no-cache && \
    apk add --no-cache openssl ca-certificates

# Create tdb user.
ENV USER=tdb
ENV UID=2001
# See https://stackoverflow.com/a/55757473/12429735RUN
RUN mkdir -p /tdb/storage
COPY telegram_dice_bot_static start_bot.sh /tdb/
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/tdb" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
        "${USER}" && chown -R "${USER}:${USER}" /tdb

USER tdb
WORKDIR /tdb

ENTRYPOINT ["sh",  "/tdb/start_bot.sh"]

