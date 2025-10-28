package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var (
	mustHearKeySegment = "must-hear"
	popularKeySegment  = "popular"
	allGenresCode      = "0"

	sortByQueryKey                  = "sort"
	userSortQueryValue              = "user"
	minReviewsQueryKey              = "reviews"
	minReviewsQueryValueAllGenres   = "500"
	minReviewsQueryValueSingleGenre = "100"
	genreQueryKey                   = "genre"

	genreCodes = []string{"15", "7", "3", "6", "132", "40", "22", "37"}

	totalScrapes = 55
)

func main() {
	_ = godotenv.Load()

	redisAddr := os.Getenv("REDIS_ADDR")
	opt, err := redis.ParseURL(redisAddr)
	if err != nil {
		log.Fatalf("Invalid REDIS_ADDR: %v", err)
	}
	rdb := redis.NewClient(opt)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("failed to connect to redis at %s: %v", redisAddr, err)
	}

	uControl := launcher.
		NewUserMode().
		Headless(true).
		Leakless(true).
		UserDataDir("tmp/t").
		Set("disable-default-apps").
		Set("no-first-run").
		Set("disable-gpu").
		NoSandbox(true).
		MustLaunch()
	browser := rod.New().ControlURL(uControl).MustConnect()
	defer browser.MustClose()

	page := browser.MustPage("https://www.albumoftheyear.org/")
	defer page.MustClose()

	scrapesDone := 0

	initialURLS := createInitialURLS()

	for filter, initialURL := range initialURLS {
		scrapeURLs := map[string]string{}

		if filter == "new" {
			scrapeURLs[filter] = initialURL
		} else {
			// starting scrape url key
			keySegments := []string{filter, popularKeySegment, allGenresCode}

			key := strings.Join(keySegments, "/")
			scrapeURLs[key] = initialURL

			u, err := url.Parse(initialURL)
			if err != nil {
				log.Println(err)
				continue
			}

			q := u.Query()
			q.Set(sortByQueryKey, userSortQueryValue)
			keySegments[1] = mustHearKeySegment
			q.Set(minReviewsQueryKey, minReviewsQueryValueAllGenres)
			u.RawQuery = q.Encode()
			key = strings.Join(keySegments, "/")
			scrapeURLs[key] = u.String()

			for i := range 2 {
				switch i {
				case 0:
					q.Set(minReviewsQueryKey, minReviewsQueryValueSingleGenre)
				case 1:
					q.Del(minReviewsQueryKey)
					q.Del(sortByQueryKey)
					keySegments[1] = popularKeySegment
				}
				for _, genre := range genreCodes {
					q.Set(genreQueryKey, genre)
					keySegments[2] = genre
					u.RawQuery = q.Encode()

					key = strings.Join(keySegments, "/")
					scrapeURLs[key] = u.String()
				}
			}

		}

		for cacheKey, scrapeURL := range scrapeURLs {
			err := scrapeAndCache(scrapeURL, cacheKey, filter, page, rdb)
			if err != nil {
				log.Println("Failed scrape and cache:", err)
				continue
			}

			scrapesDone++
			fmt.Printf("======================== %v/%v scrapes complete\n", scrapesDone, totalScrapes)

			time.Sleep(2 * time.Second)
		}
	}
}
