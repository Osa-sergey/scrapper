-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS scrapping.keywords
(
    id         SERIAL PRIMARY KEY,
    article_id INTEGER REFERENCES scrapping.articles(id) ON DELETE CASCADE,
    keyword    TEXT NOT NULL
);

CREATE INDEX idx_keywords_article_id ON scrapping.keywords(article_id);

comment on table scrapping.keywords is 'Keywords extracted from articles';
comment on column scrapping.keywords.id is 'Unique identifier';
comment on column scrapping.keywords.article_id is 'Reference to the article';
comment on column scrapping.keywords.keyword is 'Extracted keyword';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS scrapping.keywords;
-- +goose StatementEnd 