package parser

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Column struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Table struct {
	Index   int        `json:"index"`
	Columns []Column   `json:"columns"`
	Rows    [][]string `json:"rows"`
}

type Result struct {
	Tables []Table `json:"tables"`
}

func inferColumnTypes(rows [][]string) []string {
	if len(rows) == 0 {
		return nil
	}

	colCount := len(rows[0])
	types := make([]string, colCount)

	for c := 0; c < colCount; c++ {
		isInt := true
		isFloat := true

		for _, row := range rows {
			if c >= len(row) {
				continue
			}
			val := strings.TrimSpace(row[c])
			if val == "" {
				continue
			}

			if _, err := strconv.Atoi(val); err != nil {
				isInt = false
			}
			if _, err := strconv.ParseFloat(val, 64); err != nil {
				isFloat = false
			}
		}

		switch {
		case isInt:
			types[c] = "INT"
		case isFloat:
			types[c] = "FLOAT"
		default:
			types[c] = "VARCHAR"
		}
	}

	return types
}

func ParseTables(doc *goquery.Document) (*Result, error) {
	var tables []Table

	doc.Find("table").Each(func(index int, tableSel *goquery.Selection) {
		normalized, err := NormalizeTable(tableSel)
		if err != nil {
			log.Println("Error in normalizing the table")
			return
		}

		var headers []string
		var rows [][]string

		tableSel.Find("tr").First().Find("th").Each(func(_ int, s *goquery.Selection) {
			headers = append(headers, strings.TrimSpace(s.Text()))
		})

		if len(headers) == 0 {
			tableSel.Find("tr").First().Find("td").Each(func(_ int, s *goquery.Selection) {
				headers = append(headers, strings.TrimSpace(s.Text()))
			})
		}

		tableSel.Find("tr").Each(func(rowIndex int, rowSel *goquery.Selection) {
			if rowIndex == 0 {
				return
			}

			var cols []string
			rowSel.Find("td").Each(func(_ int, colSel *goquery.Selection) {
				cols = append(cols, strings.TrimSpace(colSel.Text()))
			})

			if len(cols) > 0 {
				rows = append(rows, cols)
			}
		})

		if len(headers) == 0 || len(rows) == 0 {
			log.Println("header or rows are zero while parsing")
			return
		}

		colTypes := inferColumnTypes(rows)
		var columns []Column

		for i, h := range headers {
			t := "VARCHAR"
			if i < len(colTypes) {
				t = colTypes[i]
			}
			columns = append(columns, Column{Name: h, Type: t})
		}

		tables = append(tables, Table{
			Index:   index,
			Columns: buildColumns(normalized.Headers),
			Rows:    normalized.Rows,
		})

	})

	if len(tables) == 0 {
		return nil, fmt.Errorf("no valid tables found")
	}

	return &Result{Tables: tables}, nil
}
