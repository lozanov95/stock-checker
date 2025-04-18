package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ArdesChecker struct {
	client *http.Client
}

// https://ardes.bg/product/512gb-ssd-micron-bulk-512gb-mc-nvme-2230-oem-356413
func (ac ArdesChecker) Check(url string) (CheckResponse, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return CheckResponse{}, fmt.Errorf("failed to prepare request, %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:137.0) Gecko/20100101 Firefox/137.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return CheckResponse{}, fmt.Errorf("failed to do a request, %w", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return CheckResponse{}, fmt.Errorf("failed to parse html body, %w", err)
	}

	lev := doc.Find("#price-tag").First().Text()
	cent := doc.Find(".full-price > .after-decimal").First().Text()
	price, err := strconv.ParseFloat(strings.TrimSpace(lev+cent), 64)
	if err != nil {
		return CheckResponse{}, fmt.Errorf("failed to parse price, %w", err)
	}

	available := doc.Find(".sale-action").First().Length() != 0

	return CheckResponse{
		URL:       url,
		Price:     price,
		Available: available,
	}, nil
}

func NewArdesChecker() ArdesChecker {
	return ArdesChecker{
		client: &http.Client{},
	}
}
