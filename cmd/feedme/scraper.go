package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/mahmoudk1000/feedme/internal/database"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}

func scrapeWorker(db *database.Queries, concurrency int, interval time.Duration) {
	if interval == 0 {
		interval = 60 * time.Second
	}
	log.Printf("collecting feed data every %s with concurrency %d", interval, concurrency)
	ticker := time.NewTicker(interval)

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Printf("error fetching feeds to scrape: %v", err)
			continue
		}
		log.Printf("found %d feeds to scrape", len(feeds))

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("error marking feed %d as fetched: %v", feed.ID, err)
		return
	}

	scrapedFeed, err := fetchFeed(feed.Url)
	if err != nil {
		log.Printf("error fetching feed %d from %s: %v", feed.ID, feed.Url, err)
		return
	}

	for _, item := range scrapedFeed.Channel.Item {
		if item.Title == "" {
			continue
		}
		var publishedAt sql.NullTime
		for _, layout := range []string{
			time.RFC1123, time.RFC1123Z, time.RFC822, time.RFC822Z, time.RFC3339,
		} {
			t, err := time.Parse(layout, item.PubDate)
			if err == nil {
				publishedAt = sql.NullTime{Time: t, Valid: true}
				break
			}
		}

		existingPost, err := db.CheckPostExists(context.Background(), item.Link)
		if err != nil {
			log.Printf("error checking if post exists for feed %d: %v", feed.ID, err)
			return
		}

		if !existingPost {
			_, err = db.CreatePost(context.Background(), database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
				Title:       item.Title,
				Url:         item.Link,
				PublishedAt: publishedAt,
				Description: sql.NullString{
					String: item.Description,
					Valid:  true,
				},
				FeedID: feed.ID,
			})
			if err != nil {
				log.Printf("error creating post for feed %d: %v", feed.ID, err)
			}
		}
	}
}

func fetchFeed(url string) (*RSSFeed, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("error creating request for feed %s: %v", url, err)
		return nil, err
	}
	req.Header.Set("User-Agent", "Feeder/1.0")

	client := http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := client.Do(req)

	if err != nil {
		log.Printf("error fetching data from feed %s: %v", url, err)
		return nil, err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}()

	if res.StatusCode != http.StatusOK {
		log.Printf("non-200 response from feed %s: %d", url, res.StatusCode)
		return nil, errors.New("non-200 response from feed")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("error reading body from feed %s: %v", url, err)
	}

	var feed RSSFeed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		log.Printf("error decoding XML from feed %s: %v", url, err)
		return nil, err
	}

	return &feed, nil
}
