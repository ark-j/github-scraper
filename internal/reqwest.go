package internal

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Reqwest struct {
	// http client
	client *http.Client
}

// NewReqwest return instance of Reqwest with preconfigured http.Client
func NewReqwest() *Reqwest {
	return &Reqwest{
		client: &http.Client{
			Timeout: 20 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        32,
				MaxConnsPerHost:     4,
				MaxIdleConnsPerHost: 4,
				IdleConnTimeout:     60 * time.Second,
				DisableKeepAlives:   false,
				Dial: (&net.Dialer{
					Timeout: 5 * time.Second,
				}).Dial,
				TLSHandshakeTimeout: 5 * time.Second,
			},
		},
	}
}

// Source method performs get request on url,
// parses the response body into goquery.Document
// and returns it if there is no error
func (rq *Reqwest) Source(ctx context.Context, u string) (*goquery.Document, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("not able to create request %v", err)
	}
	res, err := rq.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("not able to perform request %v", err)
	}
	defer res.Body.Close()
	return goquery.NewDocumentFromReader(res.Body)
}
