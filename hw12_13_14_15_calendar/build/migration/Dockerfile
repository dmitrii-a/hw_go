FROM golang:1.21.7-alpine3.19

ENV CODE_DIR "/code"
COPY "migrations" "${CODE_DIR}/migrations"
COPY "Makefile" ${CODE_DIR}
WORKDIR ${CODE_DIR}

RUN apk add --no-cache make
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
ENTRYPOINT make migrate-up