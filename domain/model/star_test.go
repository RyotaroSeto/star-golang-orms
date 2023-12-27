package model

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewStargazers(t *testing.T) {
	tests := []struct {
		name string
		want *Stargazers
	}{
		{
			name: "success",
			want: &Stargazers{
				Stars: make([]Stargazer, 0),
				lock:  sync.Mutex{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewStargazers())
		})
	}
}

func TestStargazers_Add(t *testing.T) {
	type fields struct {
		Stars []Stargazer
	}
	type args struct {
		stargazers []Stargazer
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantCount int
	}{
		{
			name: "success",
			fields: fields{
				Stars: make([]Stargazer, 0),
			},
			args: args{
				stargazers: []Stargazer{
					{
						StarredAt: time.Now(),
					},
				},
			},
			wantCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := &Stargazers{
				Stars: tt.fields.Stars,
				lock:  sync.Mutex{},
			}
			ss.Add(tt.args.stargazers)

			assert.Equal(t, tt.wantCount, len(ss.Stars))
		})
	}
}
