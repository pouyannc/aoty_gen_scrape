package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-rod/rod"
)

type Album struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
}

var maxAlbumsPerPage = 60

func scrapeAlbums(page *rod.Page, scrapeURLs []string) ([]Album, error) {
	nPages := len(scrapeURLs)
	totalAppended := 0
	albums := make([]Album, maxAlbumsPerPage*nPages)

	nErr := 0
	for i, u := range scrapeURLs {
		err := page.Navigate(u)
		if err != nil {
			log.Println("Page navigation error:", err)
			continue
		}

		log.Printf("Current number of pages running in browser: %v\n", len(page.Browser().MustPages()))
		log.Printf("Page html: %v\n", page.MustHTML()[:100])
		err = page.Timeout(2000*time.Millisecond).WaitElementsMoreThan(".albumBlock", 0)
		if err != nil {
			nErr++
			fmt.Println("Elements failed to load on page")
			continue
		}

		albumElements, err := page.Timeout(2000 * time.Millisecond).Elements(".albumBlock")
		if err != nil {
			fmt.Println("Failed to get album elements:", err)
			continue
		}

		appendedFromPage := 0
		for j, e := range albumElements {
			albumSlicePosition := i + j*nPages
			if albumSlicePosition >= len(albums) {
				break
			}
			albums[albumSlicePosition] = Album{
				Title:  e.MustElement(".albumTitle").MustText(),
				Artist: e.MustElement(".artistTitle").MustText(),
			}
			totalAppended++
			appendedFromPage++
			fmt.Println("appended album element text to slice:", appendedFromPage)
		}

		fmt.Println("========= appended", appendedFromPage, "albums from page")
	}
	fmt.Printf("%v/%v pages loaded elements successfully\n", len(scrapeURLs)-nErr, len(scrapeURLs))
	if nErr >= len(scrapeURLs) {
		return []Album{}, fmt.Errorf("failed to load album block elements from all pages")
	}

	fmt.Println("========= scraped ", totalAppended, "albums in total")
	return albums, nil
}
