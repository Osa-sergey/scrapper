package repository

import (
	"database/sql"
	"encoding/json"
	"strings"
)

type Article struct {
	Id          int64           `db:"id"`
	Name        string          `db:"name"`
	Text        string          `db:"text"`
	Complexity  sql.NullString  `db:"complexity"`
	ReadingTime int64           `db:"reading_time"`
	Tags        json.RawMessage `db:"tags"`
}

type ArticleInfo struct {
	ID          int             `db:"id"`
	Name        string          `db:"name"`
	Text        string          `db:"text"`
	Complexity  sql.NullString  `db:"complexity"`
	ReadingTime int             `db:"reading_time"`
	Tags        json.RawMessage `db:"tags"`
	LikeCount   int             `db:"like_count"`
	LikedByUser bool            `db:"liked_by_user"`
	Keywords    string          `db:"keywords"`
}

func (a *ArticleInfo) GetKeywords() []string {
	if a.Keywords == "" {
		return []string{}
	}
	return strings.Split(a.Keywords, ",")
}

type Keyword struct {
	ID        int64  `db:"id"`
	ArticleID int64  `db:"article_id"`
	Keyword   string `db:"keyword"`
}