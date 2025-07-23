package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"

	"github.com/mahmoudk1000/feedme/internal/database"
)

type Post struct {
	ID          uuid.UUID      `json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Title       string         `json:"title"`
	Url         string         `json:"url"`
	Description sql.NullString `json:"description"`
	PublishedAt sql.NullTime   `json:"published_at"`
	FeedID      uuid.UUID      `json:"feed_id"`
}

func DatabasePostToPost(p database.Post) Post {
	return Post{
		ID:          p.ID,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		Title:       p.Title,
		Url:         p.Url,
		Description: p.Description,
		PublishedAt: p.PublishedAt,
		FeedID:      p.FeedID,
	}
}

func DatabasePostsToPosts(posts []database.Post) []Post {
	results := make([]Post, len(posts))
	for i, post := range posts {
		results[i] = DatabasePostToPost(post)
	}
	return results
}
