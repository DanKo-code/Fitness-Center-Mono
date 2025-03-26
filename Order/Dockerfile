FROM golang:1.23.3-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o Order ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

# Устанавливаем часовой пояс GMT+3
ENV TZ=Europe/Moscow
RUN apk add --no-cache tzdata && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY --from=build /app/Order .

COPY .env .

EXPOSE ${APP_PORT}

CMD ["./Order"]