FROM golang:alpine
WORKDIR /site/

COPY go.mod ./
COPY main.go ./
COPY content ./content
COPY web ./web

RUN go mod tidy

EXPOSE 8080

RUN go build .
CMD ["./site"]

