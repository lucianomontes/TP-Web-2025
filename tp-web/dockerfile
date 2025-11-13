FROM golang

WORKDIR /app

COPY . tp-web/

WORKDIR /app/tp-web

RUN go mod download

RUN go build -o /app/web-go

EXPOSE 8081

CMD ["/app/web-go"]