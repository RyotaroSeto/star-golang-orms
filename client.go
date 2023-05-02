package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type HttpClient struct {
	url           string
	method        string
	requestHeader map[string]string
}

func (c *HttpClient) Execute() ([]byte, error) {
	res, err := c.SendRequest()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get stargazers: %s", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *HttpClient) SendRequest() (*http.Response, error) {
	req, err := http.NewRequest(c.method, c.url, nil)
	if err != nil {
		return nil, err
	}

	c.setRequestHeader(req)
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *HttpClient) setRequestHeader(req *http.Request) {
	for k, v := range c.requestHeader {
		req.Header.Set(k, v)
	}
}

func NewHttpClient(url string, method string, token string) *HttpClient {
	var hc HttpClient
	hc.requestHeader = map[string]string{"Connection": "keep-alive"}
	hc.requestHeader = map[string]string{"Authorization": "token " + token}
	hc.requestHeader = map[string]string{"Accept": "application/vnd.github.v3.star+json"}
	return &HttpClient{
		url:           url,
		method:        method,
		requestHeader: hc.requestHeader,
	}
}