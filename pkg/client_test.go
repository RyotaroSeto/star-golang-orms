package pkg

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpClient_SendRequest(t *testing.T) {
	type fields struct {
		url           string
		method        string
		requestHeader map[string]string
	}
	tests := []struct {
		name      string
		fields    fields
		want      *http.Response
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &HttpClient{
				url:           tt.fields.url,
				method:        tt.fields.method,
				requestHeader: tt.fields.requestHeader,
			}
			got, err := c.SendRequest()
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewHttpClient(t *testing.T) {
	type args struct {
		url    string
		method string
		token  string
	}
	tests := []struct {
		name string
		args args
		want *HttpClient
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewHttpClient(tt.args.url, tt.args.method, tt.args.token))
		})
	}
}

func TestHttpClient_setRequestHeader(t *testing.T) {
	type fields struct {
		url           string
		method        string
		requestHeader map[string]string
	}
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &HttpClient{
				url:           tt.fields.url,
				method:        tt.fields.method,
				requestHeader: tt.fields.requestHeader,
			}
			c.setRequestHeader(tt.args.req)
		})
	}
}
