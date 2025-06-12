import streamlit as st
import pandas as pd
import requests
from sentence_transformers import SentenceTransformer
from sklearn.metrics.pairwise import cosine_similarity
import numpy as np


@st.cache_resource
def load_embedding_model():
    return SentenceTransformer('cointegrated/rubert-tiny2')

model = load_embedding_model()

st.title("🚀 ML Engineer Roadmap")

def fetch_articles():
    query = """
    query {
      articles(page: 1, pageSize: 30) {
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
    """
    
    import os
    url = os.getenv('GRAPHQL_URL', 'http://localhost:9003/api/v1/scrapping/graph/query')
    headers = {"Content-Type": "application/json", "User-Id": "1"}
    
    try:
        response = requests.post(url, json={"query": query}, headers=headers)
        # st.json(response.json())
        if response.status_code == 200:
            return response.json()['data']['articles']['items']
        else:
            st.error(f"Error fetching articles: {response.text}")
            return []
    except Exception as e:
        st.error(f"Exception occurred: {str(e)}")
        return []

@st.cache_data
def get_article_embeddings():
    articles = fetch_articles()
    if not articles:
        return [], []
    
    texts = [article['text'] for article in articles]
    embeddings = model.encode(texts)
    return articles, embeddings

articles, article_embeddings = get_article_embeddings()

roadmap_blocks = [
    {
        "title": "1. Основы математики и статистики",
        "description": """
        - Линейная алгебра (векторы, матрицы, операции)
        - Математический анализ (производные, градиенты, оптимизация)
        - Теория вероятностей и статистика (распределения, проверка гипотез)
        - Основы статистики (среднее, дисперсия, корреляция)
        """,
        "color": "#FF6B6B"
    },
    {
        "title": "2. Основы программирования",
        "description": """
        - Программирование на Python
        - Структуры данных и алгоритмы
        - Системы контроля версий (Git)
        - Основы разработки программного обеспечения
        - Объектно-ориентированное программирование
        """,
        "color": "#4ECDC4"
    },
    {
        "title": "3. Анализ данных",
        "description": """
        - Очистка и предварительная обработка данных
        - Разведочный анализ данных (EDA)
        - Визуализация данных
        - Pandas, NumPy, Matplotlib
        - SQL и работа с базами данных
        """,
        "color": "#45B7D1"
    },
    {
        "title": "4. Основы машинного обучения",
        "description": """
        - Обучение с учителем (регрессия, классификация)
        - Обучение без учителя (кластеризация, уменьшение размерности)
        - Оценка и валидация моделей
        - Создание признаков
        - Scikit-learn
        """,
        "color": "#96CEB4"
    },
    {
        "title": "5. Глубокое обучение",
        "description": """
        - Основы нейронных сетей
        - Архитектуры CNN, RNN, LSTM
        - Трансферное обучение
        - PyTorch или TensorFlow
        - Оптимизация и настройка моделей
        """,
        "color": "#FFEEAD"
    },
    {
        "title": "6. MLOps и продакшн",
        "description": """
        - Деплой моделей
        - Контейнеризация (Docker)
        - CI/CD для ML
        - Мониторинг и поддержка моделей
        - Облачные платформы (AWS, GCP, Azure)
        """,
        "color": "#D4A5A5"
    },
    {
        "title": "7. Продвинутые темы",
        "description": """
        - Обработка естественного языка (NLP)
        - Компьютерное зрение
        - Обучение с подкреплением
        - Анализ временных рядов
        - Рекомендательные системы
        """,
        "color": "#9B59B6"
    },
    {
        "title": "8. Гибкие навыки и лучшие практики",
        "description": """
        - Решение проблем и критическое мышление
        - Коммуникация и презентация
        - Управление проектами
        - Чтение научных статей
        - Отслеживание новых тенденций
        """,
        "color": "#3498DB"
    }
]

roadmap_texts = [f"{block['title']} {block['description']}" for block in roadmap_blocks]
roadmap_embeddings = model.encode(roadmap_texts)

top_articles_indices = []
if len(article_embeddings) > 0:
    similarity_matrix = cosine_similarity(roadmap_embeddings, article_embeddings)
    top_articles_indices = np.argsort(-similarity_matrix, axis=1)[:, :3]

col1, col2 = st.columns(2)

for i, block in enumerate(roadmap_blocks):
    col = col1 if i % 2 == 0 else col2
    
    with col:
        st.markdown(
            f"""
            <div style="
                background-color: {block['color']}20;
                border-left: 5px solid {block['color']};
                padding: 20px;
                border-radius: 5px;
                margin-bottom: 20px;
            ">
                <h3 style="color: {block['color']}">{block['title']}</h3>
                <p>{block['description']}</p>
            </div>
            """,
            unsafe_allow_html=True
        )
        
        if len(articles) > 0 and len(top_articles_indices) > 0:
            st.markdown("**Recommended Articles:**")
            for idx in top_articles_indices[i]:
                article = articles[idx]
                st.markdown(
                    f"""
                    <div style="
                        background-color: {block['color']}10;
                        border-left: 3px solid {block['color']};
                        padding: 12px;
                        border-radius: 3px;
                        margin-bottom: 10px;
                    ">
                        <strong>{article['name']}</strong><br>
                        <small>Complexity: {article['complexity']} | {article['readingTime']} min read</small><br>
                        <small>Tags: {', '.join(article['tags'])}</small>
                    </div>
                    """,
                    unsafe_allow_html=True
                )
