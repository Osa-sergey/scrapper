# 🕸️ Скрапер

## 📚 Сервис для сбора и хранения статей

Система скрапинга, предоставляющая возможность собирать и хранить статьи для последующего использования и анализа.

---

## 🚀 Запуск

Чтобы запустить сервис, используйте следующую команду:

```docker compose up -d```

Это создаст и запустит все необходимые контейнеры в фоновом режиме.

---

## 🎯 Использование

### GraphQL Сервер

- Адрес сервера: http://localhost:9003/api/v1/scrapping/graph/query

  Используйте этот адрес для отправки запросов на сервер.

  Во всех запросах в Header должен быть параметр User-Id (число)

  Чтобы выключить эту проверку, установите в конфиге scraping.checkAuth = false
### GraphQL Playground

- Интерактивное окружение: http://localhost:9003/api/v1/scrapping/graph/playground

  Здесь вы можете писать и тестировать свои запросы в интерактивной среде.

```bash
query {
  articles(page: 1, pageSize: 1) {
    items {
      id
      name
      text
      complexity
      readingTime
      tags
      likes
      likedByUser
    }
    pageInfo {
      page
      pageSize
      hasNextPage
      hasPreviousPage
    }
  }
}
```

917770, 917724, 917740

```bash
query {
  article(id: 917770) {
    id
    name
    text
    complexity
    readingTime
    tags
    likes
    likedByUser
  }
}
```

```python
from sentence_transformers import SentenceTransformer
model = SentenceTransformer('cointegrated/rubert-tiny2')
model.save('local_rubert_model')
```

### REST API

- Получение статьи по ID:
  
  GET http://localhost:9003/api/v1/scrapper/article/{id}

  Эта ручка позволяет извлечь статью по её уникальному идентификатору, что может быть полезно, например, для системы рекомендаций.

- Получение списка статей по списку ID:

  POST http://localhost:9003/api/v1/scrapper/articles

  
Более подробная информация находится в [здесь](./api/api.yaml)

---