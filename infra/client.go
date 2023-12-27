package infra

import (
	"fmt"
	"net/http"
	"time"
)

const (
	responseTimeoutSec = 1000
	maxIdleTime        = 180
)

func registryHTTPClient() *http.Client {
	return &http.Client{
		Timeout: time.Duration(responseTimeoutSec) * time.Second,
		Transport: &http.Transport{
			IdleConnTimeout: time.Duration(maxIdleTime) * time.Second,
		},
	}
}

func setRequestHeader(req *http.Request) {
	for k, v := range map[string]string{
		"Connection":    "keep-alive",
		"Authorization": fmt.Sprintf("token %s", Get().GitHubToken),
		"Accept":        "application/vnd.github.v3.star+json",
	} {
		req.Header.Set(k, v)
	}
}
