package parser

type NormalizedTable struct {
	TableType string
	Headers   []string
	Rows      [][]string
	Meta      map[string]string
}
