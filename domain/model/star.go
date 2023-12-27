package model

import (
	"sync"
	"time"
)

type Stargazer struct {
	StarredAt time.Time `json:"starred_at"`
}

type Stargazers struct {
	Stars []Stargazer
	lock  sync.Mutex
}

func NewStargazers() *Stargazers {
	return &Stargazers{
		Stars: make([]Stargazer, 0),
		lock:  sync.Mutex{},
	}
}

func (ss *Stargazers) Add(stargazers []Stargazer) {
	ss.lock.Lock()
	defer ss.lock.Unlock()

	ss.Stars = append(ss.Stars, stargazers...)
}
