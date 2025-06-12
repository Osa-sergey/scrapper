import streamlit as st
import pandas as pd

st.set_page_config(
    page_title="ML Engineer Roadmap",
    page_icon="🚀",
    layout="wide"
)

st.title("🚀 ML Engineer Roadmap")

roadmap_blocks = [
    {
        "title": "1. Foundations of Mathematics & Statistics",
        "description": """
        - Linear Algebra (vectors, matrices, operations)
        - Calculus (derivatives, gradients, optimization)
        - Probability & Statistics (distributions, hypothesis testing)
        - Basic Statistics (mean, variance, correlation)
        """,
        "color": "#FF6B6B"
    },
    {
        "title": "2. Programming Fundamentals",
        "description": """
        - Python programming
        - Data structures and algorithms
        - Version control (Git)
        - Basic software engineering principles
        - Object-oriented programming
        """,
        "color": "#4ECDC4"
    },
    {
        "title": "3. Data Science & Analysis",
        "description": """
        - Data cleaning and preprocessing
        - Exploratory Data Analysis (EDA)
        - Data visualization
        - Pandas, NumPy, Matplotlib
        - SQL and database management
        """,
        "color": "#45B7D1"
    },
    {
        "title": "4. Machine Learning Basics",
        "description": """
        - Supervised Learning (regression, classification)
        - Unsupervised Learning (clustering, dimensionality reduction)
        - Model evaluation and validation
        - Feature engineering
        - Scikit-learn
        """,
        "color": "#96CEB4"
    },
    {
        "title": "5. Deep Learning",
        "description": """
        - Neural Networks fundamentals
        - CNN, RNN, LSTM architectures
        - Transfer Learning
        - PyTorch or TensorFlow
        - Model optimization and tuning
        """,
        "color": "#FFEEAD"
    },
    {
        "title": "6. MLOps & Production",
        "description": """
        - Model deployment
        - Containerization (Docker)
        - CI/CD for ML
        - Model monitoring and maintenance
        - Cloud platforms (AWS, GCP, Azure)
        """,
        "color": "#D4A5A5"
    },
    {
        "title": "7. Advanced Topics",
        "description": """
        - Natural Language Processing (NLP)
        - Computer Vision
        - Reinforcement Learning
        - Time Series Analysis
        - Recommendation Systems
        """,
        "color": "#9B59B6"
    },
    {
        "title": "8. Soft Skills & Best Practices",
        "description": """
        - Problem-solving and critical thinking
        - Communication and presentation
        - Project management
        - Research paper reading
        - Staying updated with latest trends
        """,
        "color": "#3498DB"
    }
]

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
