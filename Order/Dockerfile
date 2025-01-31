FROM golang:1.23.3-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o Order ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=build /app/Order .

COPY .env .

EXPOSE ${APP_PORT}

CMD ["./Order"]