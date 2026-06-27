FROM golang:1.26-alpine AS app

RUN apk --update --no-cache add gcc musl-dev

WORKDIR /build

COPY go.mod go.sum /build/

RUN go mod download

COPY . /build

RUN go mod tidy && go build -trimpath -o /build/stackforge ./cmd/server


FROM alpine:3.23

ENV TZ=Europe/Moscow

RUN apk --update --no-cache add tzdata curl

WORKDIR /app

COPY --from=app /build/stackforge /app/stackforge
COPY --from=app /build/config.yml /app/config.yml

ENV USER_DOCKER=stackforge
ENV UID_DOCKER=13000
ENV GID_DOCKER=13000

RUN addgroup -g ${GID_DOCKER} ${USER_DOCKER} && \
    adduser -u ${UID_DOCKER} -G ${USER_DOCKER} -s /bin/sh -D -H ${USER_DOCKER} && \
    chown -R ${USER_DOCKER}:${USER_DOCKER} /app

EXPOSE 1323/tcp

HEALTHCHECK --interval=10m --timeout=3s --start-period=5s --retries=3 CMD curl http://0.0.0.0:1323/check || exit 1

USER ${USER_DOCKER}:${USER_DOCKER}

ENTRYPOINT ["./stackforge"]
