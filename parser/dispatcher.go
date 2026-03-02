package parser

import (
	"errors"

	"github.com/PuerkitoBio/goquery"
)

var handlers = []TableHandler{
	SimpleTableHandler{},
	RowspanTableHandler{},
	KeyValueTableHandler{},
}

func NormalizeTable(table *goquery.Selection) (*NormalizedTable, error) {
	for _, h := range handlers {
		if h.CanHandle(table) {
			return h.Normalize(table)
		}
	}
	return nil, errors.New("no handler matched")
}
