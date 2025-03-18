FROM golang:1.24-bullseye

ENV DIR=/app

WORKDIR $DIR

RUN git config --global --add safe.directory $DIR &&  \
    go install github.com/air-verse/air@latest

EXPOSE 8080

CMD ["air"]
