package pkg

// type HttpClient struct {
// 	url           string
// 	method        string
// 	requestHeader map[string]string
// }

// func (c *HttpClient) SendRequest() (*http.Response, error) {
// 	req, err := http.NewRequest(c.method, c.url, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	c.setRequestHeader(req)
// 	client := new(http.Client)
// 	res, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return res, nil
// }

// func (c *HttpClient) setRequestHeader(req *http.Request) {
// 	for k, v := range c.requestHeader {
// 		req.Header.Set(k, v)
// 	}
// }

// func NewHttpClient(url string, method string, token string) *HttpClient {
// 	var hc HttpClient
// 	hc.requestHeader = map[string]string{
// 		"Connection":    "keep-alive",
// 		"Authorization": fmt.Sprintf("token %s", token),
// 		"Accept":        "application/vnd.github.v3.star+json",
// 	}
// 	return &HttpClient{
// 		url:           url,
// 		method:        method,
// 		requestHeader: hc.requestHeader,
// 	}
// }
