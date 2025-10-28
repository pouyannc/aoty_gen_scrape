package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/redis/go-redis/v9"
)

var cacheKey = "albumScrapeData"

type cachePayload struct {
	ScrapeAlbums []Album `json:"scrape_albums"`
	Ts           int64   `json:"ts"`
}

func scrapeAndCache(scrapeURL string, scrapeKey string, filter string, page *rod.Page, rdb *redis.Client) error {
	urls, err := createAllScrapeURLs(scrapeURL, filter)
	if err != nil {
		return err
	}

	scrapeData, err := scrapeAlbums(page, urls)
	if err != nil {
		return err
	}

	// cache scrapeData here
	key := fmt.Sprintf("%s:%s", cacheKey, scrapeKey)
	payload := cachePayload{
		ScrapeAlbums: scrapeData,
		Ts:           time.Now().Unix(),
	}
	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	err = rdb.Set(context.Background(), key, bytes, 0).Err()
	if err != nil {
		return err
	}

	fmt.Printf("======================== Stored new item in Redis at key: %v\n", key)
	return nil
}
