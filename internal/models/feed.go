package models

import (
	"time"

	"github.com/google/uuid"

	"github.com/mahmoudk1000/feedme/internal/database"
)

type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

func DatabaseFeedToFeed(f database.Feed) Feed {
	return Feed{
		ID:        f.ID,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
		Name:      f.Name,
		Url:       f.Url,
		UserID:    f.UserID,
	}
}

func DatabaseFeedsToFeeds(feeds []database.Feed) []Feed {
	result := make([]Feed, len(feeds))
	for i, f := range feeds {
		result[i] = DatabaseFeedToFeed(f)
	}
	return result
}
