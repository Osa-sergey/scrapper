import nltk
import ssl

try:
    _create_unverified_https_context = ssl._create_unverified_context
except AttributeError:
    pass
else:
    ssl._create_default_https_context = _create_unverified_https_context

nltk.download("stopwords")
nltk.download('punkt')

import json
import requests
from fastapi import FastAPI, HTTPException
from typing import List, Dict, Union
import nltk
from nltk.tokenize import word_tokenize
from nltk.corpus import stopwords
from sklearn.feature_extraction.text import TfidfVectorizer
import re
import os


app = FastAPI()


def preprocess_text(text: str) -> str:
    text = text.lower()
    text = re.sub(r'[^а-яА-ЯёЁ\s]', ' ', text)
    text = re.sub(r'\s+', ' ', text).strip()
    return text


def extract_keywords(text: str, n_keywords: int = 16) -> List[str]:
    text = preprocess_text(text)
    
    stop_words = set(stopwords.words('russian'))
    additional_stop_words = {'это', 'также', 'еще', 'уже', 'все', 'всего', 'всех', 'всегда', 'всегда', 'всегда'}
    stop_words.update(additional_stop_words)
    
    tokens = word_tokenize(text)
    tokens = [word for word in tokens if word not in stop_words and len(word) > 2]
    
    vectorizer = TfidfVectorizer(
        max_features=n_keywords,
        ngram_range=(1, 2),
        stop_words=list(stop_words),
        token_pattern=r'[а-яА-ЯёЁ]{3,}'
    )
    
    tfidf_matrix = vectorizer.fit_transform([text])
    feature_names = vectorizer.get_feature_names_out()
    tfidf_scores = tfidf_matrix.toarray()[0]
    
    word_scores = list(zip(feature_names, tfidf_scores))
    word_scores.sort(key=lambda x: x[1], reverse=True)
    
    return [word for word, score in word_scores[:n_keywords]]


class GraphQLClient:
    def __init__(self, base_url: str, user_id: str):
        self.base_url = base_url
        self.headers = {
            "Content-Type": "application/json",
            "User-Id": user_id
        }
    
    def get_post_text_by_id(self, post_id: int) -> str:
        query = """
        query {
            article(id: %d) {
                id
                text
            }
        }
        """ % post_id
        
        response = requests.post(
            self.base_url,
            headers=self.headers,
            json={"query": query}
        )
        
        if response.status_code == 200:
            result = response.json()
            article = result.get('data', {}).get('article', {})
            if article:
                return article.get('text', '')
            return "Article not found"
        else:
            return f"Error: {response.status_code} - {response.text}"


client = GraphQLClient(
    os.getenv('GRAPHQL_URL', 'http://localhost:9003/api/v1/scrapping/graph/query'),
    os.getenv('USER_ID', '1')
)


@app.get("/extract_keywords_by_id/{post_id}", response_model=Dict[str, Union[List[str], str]])
async def extract_keywords_by_id(post_id: int) -> dict:
    text = client.get_post_text_by_id(post_id)
    if text.startswith("Error:") or text == "Article not found":
        raise HTTPException(status_code=404, detail=text)
    
    keywords = extract_keywords(text)
    return {
        "keywords": keywords,
        "text": text
    }


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
