package parser

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type TableHandler interface {
	CanHandle(table *goquery.Selection) bool
	Normalize(table *goquery.Selection) (*NormalizedTable, error)
}

func extractHeaders(table *goquery.Selection) []string {
	var headers []string

	// 1️⃣ Prefer <th> headers (first header row)
	table.Find("tr").EachWithBreak(func(i int, tr *goquery.Selection) bool {
		th := tr.Find("th")
		if th.Length() > 0 {
			th.Each(func(_ int, h *goquery.Selection) {
				text := strings.TrimSpace(h.Text())
				headers = append(headers, text)
			})
			return false // stop after first <th> row
		}
		return true
	})

	// 2️⃣ Fallback: first <td> row if no <th>
	if len(headers) == 0 {
		table.Find("tr").First().Find("td").Each(func(_ int, td *goquery.Selection) {
			text := strings.TrimSpace(td.Text())
			headers = append(headers, text)
		})
	}

	// 3️⃣ Absolute fallback: generate placeholders
	if len(headers) == 0 {
		colCount := table.Find("tr").First().Find("td").Length()
		for i := 0; i < colCount; i++ {
			headers = append(headers, "")
		}
	}

	return headers
}
