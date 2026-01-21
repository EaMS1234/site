FROM golang:alpine
WORKDIR /site/

ENV TZ="America/Sao_Paulo"
RUN date

copy . .

RUN go mod tidy

EXPOSE 8080

RUN go build .

USER 1234:1234

CMD ["./site"]

