package models

import (
	"time"

	"github.com/google/uuid"

	"github.com/mahmoudk1000/feedme/internal/database"
)

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FeedID    uuid.UUID `json:"feed_id"`
	UserID    uuid.UUID `json:"user_id"`
}

func DatabaseFeedFollowToFeedFollow(f database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        f.ID,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
		FeedID:    f.FeedID,
		UserID:    f.UserID,
	}
}

func DatabaseFeedFollowsToFeedFollows(feeds []database.FeedFollow) []FeedFollow {
	feedFollows := make([]FeedFollow, len(feeds))
	for i, f := range feeds {
		feedFollows[i] = DatabaseFeedFollowToFeedFollow(f)
	}
	return feedFollows
}
