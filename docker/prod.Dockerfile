ARG DIR=/app

FROM golang:1.24-bullseye AS build

ARG DIR
WORKDIR $DIR

COPY . .

RUN git config --global --add safe.directory $DIR &&  \
    go build

FROM debian:bullseye-slim

ARG DIR
WORKDIR $DIR

COPY --from=build ${DIR}/challenge-bravo .

EXPOSE 8080

CMD ["./challenge-bravo"]
