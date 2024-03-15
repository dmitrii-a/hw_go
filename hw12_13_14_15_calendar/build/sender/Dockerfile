FROM golang:1.21.7-alpine3.19 as build

ENV CODE_DIR "/code"
ENV BIN_SENDER "${CODE_DIR}/sender"
ENV CGO_ENABLED 0
ARG LDFLAGS

COPY . ${CODE_DIR}
WORKDIR ${CODE_DIR}

RUN apk add --no-cache make
RUN go mod download

RUN make build-sender

FROM alpine:3.19.1
ENV CODE_DIR "/code"
ENV BIN_SENDER "${CODE_DIR}/sender"

COPY --from=build ${BIN_SENDER} ${BIN_SENDER}

ENV CONFIG_FILE "${CODE_DIR}/configs/config.yaml"
COPY ./configs/config.yaml ${CONFIG_FILE}

ENTRYPOINT ${BIN_SENDER} -config ${CONFIG_FILE}