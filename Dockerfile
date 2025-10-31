FROM golang:alpine
WORKDIR /site/

copy . .

RUN go mod tidy

EXPOSE 8080

RUN go build .
CMD ["./site"]

