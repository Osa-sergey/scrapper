package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var errNotFound = errors.New("obj not found")

func IsNotFoundError(err error) bool {
	return errors.Is(err, errNotFound)
}

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetLastArticle(ctx context.Context) (int64, error) {
	var id sql.NullInt64
	err := r.db.GetContext(ctx, &id, "SELECT MAX(id) from scrapping.articles")
	if err != nil {
		return 0, fmt.Errorf("error in db: %v", err)
	}
	if !id.Valid {
		return 0, errNotFound
	}
	return id.Int64, nil
}

func (r *Repository) GetFirstArticle(ctx context.Context) (int64, error) {
	var id sql.NullInt64
	err := r.db.GetContext(ctx, &id, "SELECT MIN(id) from scrapping.articles")
	if err != nil {
		return 0, fmt.Errorf("error in db: %v", err)
	}
	if !id.Valid {
		return 0, errNotFound
	}
	return id.Int64, nil
}

func (r *Repository) AddArticle(ctx context.Context, article *Article) error {
	stmt, err := r.db.PrepareNamedContext(ctx, `INSERT INTO scrapping.articles 
    (id, name, text, complexity, reading_time, tags)
	VALUES (:id, :name, :text, :complexity, :reading_time, :tags)
	RETURNING id`)

	if err != nil {
		return fmt.Errorf("add article error: %v", err)
	}
	defer stmt.Close()

	err = stmt.GetContext(ctx, &article.Id, article)
	if err != nil {
		return fmt.Errorf("add article, error in GetContext: %v", err)
	}
	return nil
}

func (r *Repository) GetArticleById(ctx context.Context, id int) (*Article, error) {
	article := &Article{}
	err := r.db.GetContext(ctx, article, "SELECT * FROM scrapping.articles WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errNotFound
		}
		return nil, fmt.Errorf("GetArticleById db error: %v", err)
	}
	return article, nil
}

func (r *Repository) GetArticlesInfo(ctx context.Context, userId int, cursor *Cursor) ([]*ArticleInfo, *PaginationInfo, error) {
	if cursor != nil {
		err := cursor.Validate()
		if err != nil {
			return nil, nil, fmt.Errorf("%w: %v", ErrInvalidCursor, err)
		}
	} else {
		return nil, nil, fmt.Errorf("cursor is empty")
	}

	query := `
		WITH liked_articles AS (
			SELECT article_id
			FROM scrapping.likes
			WHERE user_id = $1
		),
		article_likes AS (
			SELECT article_id, COUNT(*) AS like_count
			FROM scrapping.likes
			GROUP BY article_id
		),
		article_keywords AS (
			SELECT article_id, array_to_string(array_agg(keyword), ',') as keywords
			FROM scrapping.keywords
			GROUP BY article_id
		)
		SELECT a.id, a.name, a.text, a.complexity, a.reading_time, a.tags,
			   COALESCE(al.like_count, 0) AS like_count,
			   CASE WHEN la.article_id IS NOT NULL THEN true ELSE false END AS liked_by_user,
			   COALESCE(ak.keywords, '') as keywords
		FROM scrapping.articles a
		LEFT JOIN article_likes al ON a.id = al.article_id
		LEFT JOIN liked_articles la ON a.id = la.article_id
		LEFT JOIN article_keywords ak ON a.id = ak.article_id
		ORDER BY a.id
		%s;
    `
	limit := cursor.limitPlusOne()

	query = fmt.Sprintf(query, limit)

	result := make([]*ArticleInfo, 0)

	err := r.db.SelectContext(ctx, &result, query, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, errNotFound
		}
		return nil, nil, fmt.Errorf("error in db: %v", err)
	}

	hasNextPage := false
	if len(result) > cursor.limit() {
		hasNextPage = true
		result = result[:len(result)-1]
	}
	return result, &PaginationInfo{
		HasNextPage:     hasNextPage,
		HasPreviousPage: cursor.offset() > 0,
	}, nil
}

func (r *Repository) GetArticleInfoById(ctx context.Context, userId, articleId int) (*ArticleInfo, error) {
	query := `
		WITH liked_articles AS (
			SELECT article_id
			FROM scrapping.likes
			WHERE user_id = $1
		),
		article_likes AS (
			SELECT article_id, COUNT(*) AS like_count
			FROM scrapping.likes
			GROUP BY article_id
		),
		article_keywords AS (
			SELECT article_id, array_to_string(array_agg(keyword), ',') as keywords
			FROM scrapping.keywords
			GROUP BY article_id
		)
		SELECT a.id, a.name, a.text, a.complexity, a.reading_time, a.tags,
			   COALESCE(al.like_count, 0) AS like_count,
			   CASE WHEN la.article_id IS NOT NULL THEN true ELSE false END AS liked_by_user,
			   COALESCE(ak.keywords, '') as keywords
		FROM scrapping.articles a
		LEFT JOIN article_likes al ON a.id = al.article_id
		LEFT JOIN liked_articles la ON a.id = la.article_id
		LEFT JOIN article_keywords ak ON a.id = ak.article_id
		WHERE id = $2
		ORDER BY a.id;
    `

	var article ArticleInfo
	err := r.db.GetContext(ctx, &article, query, userId, articleId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errNotFound
		}
		return nil, fmt.Errorf("db error: %v", err)
	}

	return &article, nil
}

func (r *Repository) Like(ctx context.Context, userId, articleId int) error {
	query := `INSERT INTO scrapping.likes (user_id, article_id) VALUES ($1, $2)`

	_, err := r.db.ExecContext(ctx, query, userId, articleId)
	if err != nil {
		var pqErr *pgconn.PgError
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23503":
				return fmt.Errorf("non exist article")
			case "23505":
				return fmt.Errorf("like already exist")
			}
		}
		return fmt.Errorf("error in db: %v", err)
	}
	return nil
}

func (r *Repository) Unlike(ctx context.Context, userId, articleId int) error {
	query := `DELETE FROM scrapping.likes WHERE user_id = $1 AND article_id = $2`
	rows, err := r.db.ExecContext(ctx, query, userId, articleId)
	if err != nil {
		return err
	}
	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no likes were found")
	}
	return nil
}

func (r *Repository) GetArticlesByIds(ctx context.Context, userId int, ids []int) ([]*ArticleInfo, error) {
	query := `
		WITH liked_articles AS (
			SELECT article_id
			FROM scrapping.likes
			WHERE user_id = $1
		),
		article_likes AS (
			SELECT article_id, COUNT(*) AS like_count
			FROM scrapping.likes
			GROUP BY article_id
		),
		article_keywords AS (
			SELECT article_id, array_to_string(array_agg(keyword), ',') as keywords
			FROM scrapping.keywords
			GROUP BY article_id
		)
		SELECT a.id, a.name, a.text, a.complexity, a.reading_time, a.tags,
			   COALESCE(al.like_count, 0) AS like_count,
			   CASE WHEN la.article_id IS NOT NULL THEN true ELSE false END AS liked_by_user,
			   COALESCE(ak.keywords, '') as keywords
		FROM scrapping.articles a
		LEFT JOIN article_likes al ON a.id = al.article_id
		LEFT JOIN liked_articles la ON a.id = la.article_id
		LEFT JOIN article_keywords ak ON a.id = ak.article_id
		WHERE a.id = ANY($2)
		ORDER BY a.id;
    `

	var articles []*ArticleInfo
	err := r.db.SelectContext(ctx, &articles, query, userId, pq.Array(ids))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errNotFound
		}
		return nil, fmt.Errorf("db error: %v", err)
	}

	return articles, nil
}

func (r *Repository) Search(ctx context.Context, userId int, search string, pageSize int) ([]*ArticleInfo, error) {
	query := `
		WITH liked_articles AS (
			SELECT article_id
			FROM scrapping.likes
			WHERE user_id = $1
		),
		article_likes AS (
			SELECT article_id, COUNT(*) AS like_count
			FROM scrapping.likes
			GROUP BY article_id
		),
		article_keywords AS (
			SELECT article_id, array_to_string(array_agg(keyword), ',') as keywords
			FROM scrapping.keywords
			GROUP BY article_id
		)
		SELECT a.id, a.name, a.text, a.complexity, a.reading_time, a.tags,
			   COALESCE(al.like_count, 0) AS like_count,
			   CASE WHEN la.article_id IS NOT NULL THEN true ELSE false END AS liked_by_user,
			   COALESCE(ak.keywords, '') as keywords
		FROM scrapping.articles a
		LEFT JOIN article_likes al ON a.id = al.article_id
		LEFT JOIN liked_articles la ON a.id = la.article_id
		LEFT JOIN article_keywords ak ON a.id = ak.article_id
		WHERE a.name ILIKE $2 OR a.text ILIKE $2
		ORDER BY a.id
		LIMIT $3;
    `

	var articles []*ArticleInfo
	err := r.db.SelectContext(ctx, &articles, query, userId, "%"+search+"%", pageSize)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errNotFound
		}
		return nil, fmt.Errorf("db error: %v", err)
	}

	return articles, nil
}

func (r *Repository) AddKeywords(ctx context.Context, articleId int64, keywords []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction error: %v", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `INSERT INTO scrapping.keywords 
		(article_id, keyword)
		VALUES ($1, $2)`)
	if err != nil {
		return fmt.Errorf("prepare statement error: %v", err)
	}
	defer stmt.Close()

	for _, keyword := range keywords {
		_, err = stmt.ExecContext(ctx, articleId, keyword)
		if err != nil {
			return fmt.Errorf("insert keyword error: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit transaction error: %v", err)
	}

	return nil
}

func (r *Repository) GetKeywordsByArticleId(ctx context.Context, articleId int64) ([]string, error) {
	query := `
		SELECT keyword
		FROM scrapping.keywords
		WHERE article_id = $1
		ORDER BY keyword;
	`

	var keywords []string
	err := r.db.SelectContext(ctx, &keywords, query, articleId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errNotFound
		}
		return nil, fmt.Errorf("db error: %v", err)
	}
	return keywords, nil
}
