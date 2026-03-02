package parser

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type RowspanTableHandler struct{}

type spanCell struct {
	value string
	left  int
}

func (h RowspanTableHandler) CanHandle(t *goquery.Selection) bool {
	found := false
	t.Find("td").Each(func(_ int, s *goquery.Selection) {
		if _, ok := s.Attr("rowspan"); ok {
			found = true
		}
	})
	return found
}

func (h RowspanTableHandler) Normalize(t *goquery.Selection) (*NormalizedTable, error) {
	headers := extractHeaders(t)
	colCount := len(headers)

	type spanCell struct {
		value string
		left  int
	}

	active := map[int]*spanCell{}
	var rows [][]string

	t.Find("tr").Each(func(i int, tr *goquery.Selection) {
		if i == 0 {
			return
		}

		row := make([]string, colCount)
		col := 0

		// 1️⃣ Fill active rowspans
		for col < colCount {
			if sp, ok := active[col]; ok {
				row[col] = sp.value
				sp.left--
				if sp.left == 0 {
					delete(active, col)
				}
				col++
			} else {
				break
			}
		}

		// 2️⃣ Process actual <td> cells
		tr.Find("td").Each(func(_ int, td *goquery.Selection) {

			// Advance to next free column
			for col < colCount {
				if _, ok := active[col]; !ok && row[col] == "" {
					break
				}
				col++
			}

			// 🚨 HARD STOP: ignore overflow cells
			if col >= colCount {
				return
			}

			val := strings.TrimSpace(td.Text())
			row[col] = val

			if rs, ok := td.Attr("rowspan"); ok {
				if n, err := strconv.Atoi(rs); err == nil && n > 1 {
					active[col] = &spanCell{
						value: val,
						left:  n - 1,
					}
				}
			}

			col++
		})

		rows = append(rows, row)
	})

	return &NormalizedTable{
		TableType: "rowspan",
		Headers:   headers,
		Rows:      rows,
	}, nil
}
