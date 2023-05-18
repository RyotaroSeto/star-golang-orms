package main

import "regexp"

func getStarPageURL(headerLink string) (string, string) {
	var nextPage, lastPage string
	nextRe := regexp.MustCompile("<([^>]+)>; rel=\"next\"")
	nextMatch := nextRe.FindStringSubmatch(headerLink)
	nextPage = nextMatch[1]

	lastRe := regexp.MustCompile("<([^>]+)>; rel=\"last\"")
	lastMatch := lastRe.FindStringSubmatch(headerLink)
	lastPage = lastMatch[1]

	return nextPage, lastPage
}
