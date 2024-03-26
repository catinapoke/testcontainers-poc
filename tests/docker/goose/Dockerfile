FROM golang:bullseye

RUN go install https://github.com/pressly/goose
ADD migrations .

ENTRYPOINT ["goose"]