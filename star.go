package main

import "regexp"

func getNextPageURL(headerLink string) string {
	var nextPage string
	nextRe := regexp.MustCompile("<([^>]+)>; rel=\"next\"")
	nextMatch := nextRe.FindStringSubmatch(headerLink)
	nextPage = nextMatch[1]

	return nextPage
}

func getLastPageURL(headerLink string) string {
	var lastPage string
	lastRe := regexp.MustCompile("<([^>]+)>; rel=\"last\"")
	lastMatch := lastRe.FindStringSubmatch(headerLink)
	lastPage = lastMatch[1]

	return lastPage
}
