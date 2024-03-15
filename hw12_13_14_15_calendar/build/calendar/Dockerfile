FROM golang:1.21.7-alpine3.19 as build

ENV CODE_DIR "/code"
ENV BIN_CALENDAR "${CODE_DIR}/calendar"
ENV CGO_ENABLED 0
ARG LDFLAGS

COPY . ${CODE_DIR}
WORKDIR ${CODE_DIR}

RUN apk add --no-cache make
RUN go mod download

RUN make build-calendar

FROM alpine:3.19.1
ENV CODE_DIR "/code"
ENV BIN_CALENDAR "${CODE_DIR}/calendar"

COPY --from=build ${BIN_CALENDAR} ${BIN_CALENDAR}

ENV CONFIG_FILE "${CODE_DIR}/configs/config.yaml"
COPY ./configs/config.yaml ${CONFIG_FILE}

ENTRYPOINT ${BIN_CALENDAR} -config ${CONFIG_FILE}