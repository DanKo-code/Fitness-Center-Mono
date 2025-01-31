# Dockerfile.Frontend

# Устанавливаем Node.js
FROM node:18-alpine AS build

# Устанавливаем рабочую директорию
WORKDIR /usr/src/app

# Копируем package.json и package-lock.json
COPY package*.json ./

# Устанавливаем зависимости
RUN npm install

# Копируем все файлы в контейнер
COPY . .

# Сборка приложения
RUN npm run build

# Сервируем приложение с помощью сервера Nginx
FROM nginx:alpine

# Копируем package.json и package-lock.json
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Копируем сгенерированные файлы из предыдущего этапа сборки
COPY --from=build /usr/src/app/build /usr/share/nginx/html

# Экспонируем порт для доступа к приложению
EXPOSE 80

# Запускаем Nginx
CMD ["nginx", "-g", "daemon off;"]