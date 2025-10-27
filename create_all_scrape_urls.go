package main

import (
	"net/url"
	"strconv"
	"strings"
)

var nYears = 2
var nMonths = 3
var pathMonths = []string{
	"january-01.php",
	"february-02.php",
	"march-03.php",
	"april-04.php",
	"may-05.php",
	"june-06.php",
	"july-07.php",
	"august-08.php",
	"september-09.php",
	"october-10.php",
	"november-11.php",
	"december-12.php",
}

func createAllScrapeURLs(initialURL, filter string) ([]string, error) {
	resSlice := []string{
		initialURL,
	}

	initialParsed, err := url.Parse(initialURL)
	if err != nil {
		return []string{}, err
	}
	pathSegments := strings.Split(strings.Trim(initialParsed.Path, "/"), "/")

	switch filter {
	case "months":
		for range nMonths {
			currPathMonth := pathSegments[2]
			// path for months always ends in '.php'
			currIndex, err := strconv.Atoi(currPathMonth[len(currPathMonth)-6 : len(currPathMonth)-4])
			if err != nil {
				return []string{}, err
			}
			prevIndex := currIndex - 2

			if prevIndex < 0 {
				prevIndex = 11
				pathSegments, err = decreasePathYear(pathSegments)
				if err != nil {
					return []string{}, err
				}
			}

			pathSegments[2] = pathMonths[prevIndex]
			initialParsed.Path = "/" + strings.Join(pathSegments, "/")
			resSlice = append(resSlice, initialParsed.String())
		}
	case "years":
		for range nYears {
			pathSegments, err = decreasePathYear(pathSegments)
			if err != nil {
				return []string{}, err
			}

			// the years urls always end the path in an extra '/'
			initialParsed.Path = "/" + strings.Join(pathSegments, "/") + "/"
			resSlice = append(resSlice, initialParsed.String())
		}
	}

	return resSlice, nil
}

func decreasePathYear(segments []string) ([]string, error) {
	year, err := strconv.Atoi(segments[0])
	if err != nil {
		return []string{}, err
	}
	segments[0] = strconv.Itoa(year - 1)
	return segments, nil
}
