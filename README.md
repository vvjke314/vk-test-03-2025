# VK Test Project

Проект представляет собой сервис для работы с хранилищем ключ-значение, использующий Tarantool в качестве базы данных.

## Требования

- Docker
- Docker Compose
- Git

## Структура проекта

```
.
├── cmd/                    # Точка входа приложения
├── config/                 # Конфигурационные файлы
├── internal/              # Внутренние пакеты
├── pkg/                   # Публичные пакеты
├── tarantool/            # Конфигурация и скрипты Tarantool
├── .env                  # Файл с переменными окружения
├── docker-compose.yaml   # Конфигурация Docker Compose
└── Dockerfile           # Конфигурация сборки приложения
```

## Настройка окружения

1. Создайте файл `.env` в корневой директории проекта со следующим содержимым:

```env
TRNTLUSER=admin
TRNTLPASS=admin
TRNTLPORT=3301
```

## Запуск проекта

### Локальный запуск

1. Клонируйте репозиторий:
```bash
git clone <https://github.com/vvjke314/vk-test-03-2025.git>
cd vk-test-03-2025
```

2. Соберите и запустите контейнеры:
```bash
docker-compose up --build
```

После успешного запуска:
- Приложение будет доступно по адресу: `http://localhost:8080`
- Tarantool будет доступен по порту: `3301`

### Остановка проекта

Для остановки всех контейнеров выполните:
```bash
docker-compose down
```

## Деплой на сервер

1. Скопируйте все файлы проекта на сервер:
   - `docker-compose.yaml`
   - `Dockerfile`
   - `.env`
   - Исходный код проекта

2. Убедитесь, что на сервере установлены Docker и Docker Compose

3. Создайте файл `.env` на сервере с необходимыми переменными окружения

4. Запустите проект:
```bash
docker-compose up -d --build
```

Флаг `-d` запустит контейнеры в фоновом режиме.

### Мониторинг логов

Для просмотра логов приложения:
```bash
docker-compose logs -f app
```

Для просмотра логов Tarantool:
```bash
docker-compose logs -f tarantool
```

### Перезапуск сервисов

При необходимости перезапуска отдельных сервисов:
```bash
docker-compose restart app      # Перезапуск приложения
docker-compose restart tarantool  # Перезапуск Tarantool
```

## Устранение неполадок

1. Если приложение не может подключиться к Tarantool:
   - Проверьте, что все переменные окружения в `.env` установлены корректно
   - Убедитесь, что порты не заняты другими сервисами
   - Проверьте логи обоих сервисов

2. Если контейнеры не запускаются:
   - Проверьте логи: `docker-compose logs`
   - Убедитесь, что все необходимые порты свободны
   - Проверьте права доступа к директориям с данными

## Очистка данных

Для полной очистки данных и пересоздания контейнеров:
```bash
docker-compose down -v
docker-compose up --build
```

⚠️ **Внимание**: Эта команда удалит все данные в Tarantool!
