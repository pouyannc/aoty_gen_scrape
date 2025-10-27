package main

import (
	"context"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var genreCodes = []string{"15", "7", "3", "6", "132", "40", "22", "37"}

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

	initialURLS := createInitialURLS()

	for filter, initialURL := range initialURLS {
		scrapeURLs := []string{}

		if filter == "new" {
			scrapeURLs = append(scrapeURLs, initialURL)
		} else {
			scrapeURLs = append(scrapeURLs, initialURL)

			u, err := url.Parse(initialURL)
			if err != nil {
				log.Println(err)
				continue
			}

			q := u.Query()
			q.Set("sort", "user")
			q.Set("reviews", "500")
			u.RawQuery = q.Encode()
			scrapeURLs = append(scrapeURLs, u.String())

			for i := range 2 {
				switch i {
				case 0:
					q.Set("reviews", "100")
				case 1:
					q.Del("reviews")
					q.Del("sort")
				}
				for _, genre := range genreCodes {
					q.Set("genre", genre)
					u.RawQuery = q.Encode()

					scrapeURLs = append(scrapeURLs, u.String())
				}
			}

		}

		for _, scrapeURL := range scrapeURLs {
			err := scrapeAndCache(scrapeURL, filter, page, rdb)
			if err != nil {
				log.Println("Failed scrape and cache:", err)
				continue
			}

			time.Sleep(2 * time.Second)
		}
	}
}
