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
	client *http.Client
}

func NewReqwest() *Reqwest {
	return &Reqwest{
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxConnsPerHost:     5,
				MaxIdleConnsPerHost: 5,
				IdleConnTimeout:     30 * time.Second,
				DisableKeepAlives:   false,
				Dial: (&net.Dialer{
					Timeout: 5 * time.Second,
				}).Dial,
				TLSHandshakeTimeout: 5 * time.Second,
			},
		},
	}
}

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
