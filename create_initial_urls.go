package main

import (
	"fmt"
	"time"
)

var (
	baseURL = "https://www.albumoftheyear.org"
)

func createInitialURLS() map[string]string {
	year, month, _ := time.Now().Date()
	monthNum := int(month)

	initialURLS := map[string]string{}

	initialURLS["new"] = baseURL + "/releases/this-week/"
	initialURLS["months"] = fmt.Sprintf("%s/%d/releases/%s?type=lp", baseURL, year, pathMonths[monthNum-1])
	initialURLS["year"] = fmt.Sprintf("%s/%d/releases/?type=lp", baseURL, year)
	initialURLS["years"] = fmt.Sprintf("%s/%d/releases/?type=lp", baseURL, year)

	return initialURLS
}
