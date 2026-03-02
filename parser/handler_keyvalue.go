package parser

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type KeyValueTableHandler struct{}

func (h KeyValueTableHandler) CanHandle(t *goquery.Selection) bool {
	first := t.Find("tr").First()
	return first.Find("th").Length() == 0 &&
		first.Find("td").Length() == 2
}

func (h KeyValueTableHandler) Normalize(t *goquery.Selection) (*NormalizedTable, error) {
	headers := []string{"key", "value"}
	var rows [][]string

	t.Find("tr").Each(func(_ int, tr *goquery.Selection) {
		tds := tr.Find("td")
		if tds.Length() != 2 {
			return
		}
		key := strings.TrimSpace(tds.First().Text())
		val := strings.TrimSpace(tds.Last().Text())
		rows = append(rows, []string{key, val})
	})

	return &NormalizedTable{
		TableType: "key_value",
		Headers:   headers,
		Rows:      rows,
	}, nil
}
