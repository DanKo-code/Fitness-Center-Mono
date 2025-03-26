# Устанавливаем Node.js и билдим фронтенд
FROM node:18-alpine AS build

WORKDIR /usr/src/app

COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

# Устанавливаем Nginx
FROM nginx:alpine

# Копируем конфиг Nginx
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Копируем собранные файлы React
COPY --from=build /usr/src/app/build /usr/share/nginx/html

# Копируем SSL-сертификаты
COPY cert/cert.pem /etc/nginx/ssl/cert.pem
COPY cert/key.pem /etc/nginx/ssl/key.pem

# Открываем порты HTTPS и HTTP
EXPOSE 81 443

CMD ["nginx", "-g", "daemon off;"]
