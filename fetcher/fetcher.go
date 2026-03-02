package fetcher

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func FetchDocument(url string) (*goquery.Document, error) {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 🔴 REQUIRED for Wikipedia
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; TableCollector/1.0)")
	req.Header.Set("Accept", "text/html")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var reader io.Reader = resp.Body

	// 🔴 Handle gzip
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		defer gz.Close()
		reader = gz
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status %d", resp.StatusCode)
	}

	return goquery.NewDocumentFromReader(reader)
}
