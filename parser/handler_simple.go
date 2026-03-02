package parser

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type SimpleTableHandler struct{}

func (h SimpleTableHandler) CanHandle(t *goquery.Selection) bool {
	hasSpan := false
	t.Find("td,th").Each(func(_ int, s *goquery.Selection) {
		if _, ok := s.Attr("rowspan"); ok {
			hasSpan = true
		}
		if _, ok := s.Attr("colspan"); ok {
			hasSpan = true
		}
	})
	return !hasSpan
}

func (h SimpleTableHandler) Normalize(t *goquery.Selection) (*NormalizedTable, error) {
	headers := extractHeaders(t)
	var rows [][]string

	t.Find("tr").Each(func(i int, tr *goquery.Selection) {
		if i == 0 {
			return
		}
		var row []string
		tr.Find("td").Each(func(_ int, td *goquery.Selection) {
			row = append(row, strings.TrimSpace(td.Text()))
		})
		if len(row) > 0 {
			rows = append(rows, row)
		}
	})

	return &NormalizedTable{
		TableType: "simple",
		Headers:   headers,
		Rows:      rows,
	}, nil
}
