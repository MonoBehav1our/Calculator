# frontend/Dockerfile
FROM node:18-alpine

WORKDIR /app

# Копируем зависимости и устанавливаем их
COPY package*.json ./
RUN npm install

# Копируем весь исходный код
COPY . .

EXPOSE 3000

# Команда для локальной разработки (npm start)
CMD ["npm", "start"]