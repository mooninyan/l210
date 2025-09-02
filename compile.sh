#!/bin/bash

# Название выходного файла
OUTPUT_NAME="build/l210"

# Создаем директорию build, если её нет
mkdir -p build

# Установка переменных окружения для Linux x86_64
export GOOS=linux
export GOARCH=amd64

# Выполнение сборки
echo "Сборка проекта для Linux..."
go build -o "$OUTPUT_NAME"

if [ $? -eq 0 ]; then
    echo "Сборка завершена успешно: $OUTPUT_NAME"
else
    echo "Ошибка при сборке"
    exit 1
fi