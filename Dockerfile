# Используем Buffalo-образ для сборки
FROM gobuffalo/buffalo:v0.18.14 as builder

# Установка Node.js и Yarn через официальный Node.js образ
RUN curl -sL https://deb.nodesource.com/setup_16.x | bash - \
    && apt-get install -y nodejs yarn

# Создаём директорию для приложения
RUN mkdir -p /src/sharaphka
WORKDIR /src/sharaphka

# Добавляем package.json и yarn.lock перед установкой зависимостей
ADD package.json .
ADD yarn.lock .

# Устанавливаем зависимости для фронтенда
RUN yarn install

# Копируем Go-зависимости и устанавливаем их
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
ADD . .

# Собираем Buffalo приложение
RUN buffalo build --static -o /bin/app

# Используем финальный минимальный образ на основе Alpine
FROM alpine
RUN apk add --no-cache bash ca-certificates

WORKDIR /bin/
COPY --from=builder /bin/app .

# Запускаем приложение
CMD exec /bin/app
