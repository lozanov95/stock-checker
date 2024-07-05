package main

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type OzoneChecker struct {
	Replacer *strings.Replacer
}

func (oc *OzoneChecker) isItemAvailable(doc *goquery.Document) bool {
	return doc.Find("#availability-holder.availability.out-of-stock").Length() == 0
}

func (oc *OzoneChecker) getPrice(doc *goquery.Document) float64 {
	priceNode := doc.Find(".product-options [id^='product-price-']")
	if priceNode.Length() == 0 {
		return 0
	}

	num, err := strconv.ParseFloat(
		strings.TrimSpace(
			oc.Replacer.Replace(
				priceNode.Text(),
			),
		), 64)

	if err != nil {
		slog.Error(err.Error())
		return 0
	}

	return num
}

func (oc *OzoneChecker) Check(url string) (CheckResponse, error) {
	response := CheckResponse{URL: url}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return response, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return response, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return response, err
	}

	response.Available = oc.isItemAvailable(doc)
	response.Price = oc.getPrice(doc)

	return response, nil
}

func NewOzoneChecker() *OzoneChecker {
	return &OzoneChecker{
		Replacer: strings.NewReplacer(
			",", ".",
			"лв.", "",
			" ", ""),
	}
}
