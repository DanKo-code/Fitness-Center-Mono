FROM alpine:latest

RUN apk add --no-cache curl unzip

RUN curl -fsSL https://bin.equinox.io/c/bNyj1mQVY4c/ngrok-stable-linux-amd64.zip -o ngrok.zip && \
    unzip ngrok.zip && \
    mv ngrok /usr/local/bin/ngrok && \
    rm ngrok.zip

WORKDIR /app

EXPOSE 4040

ARG NGROK_AUTHTOKEN
RUN ngrok config add-authtoken ${NGROK_AUTHTOKEN}

ENTRYPOINT ["ngrok"]
CMD ["http", "80"]
