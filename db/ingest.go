package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
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

func InsertRows(db *sql.DB, tableName string, columns []Column, rows [][]string) error {
	var colNames []string
	for i, c := range columns {
		name := SanitizeColumnName(c.Name, i)
		if name == "" {
			continue
		}
		colNames = append(colNames, "`"+name+"`")
	}

	if len(colNames) == 0 {
		log.Println("Column legth is zero: ", tableName)
		return nil
	}

	placeholder := "(" + strings.TrimRight(strings.Repeat("?,", len(colNames)), ",") + ")"
	var placeholders []string
	var values []any

	for _, row := range rows {
		if len(row) < len(colNames) {
			continue
		}

		placeholders = append(placeholders, placeholder)
		for i := 0; i < len(colNames); i++ {
			values = append(values, row[i])
		}
	}
	log.Println(tableName)
	// 🔴 THIS IS THE CRITICAL GUARD
	if len(placeholders) == 0 {
		log.Println("no valid rows to insert in: ", tableName)
		return nil // no valid rows to insert
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES %s",
		tableName,
		strings.Join(colNames, ","),
		strings.Join(placeholders, ","),
	)

	log.Println("TABLE INSERT QUERY:")
	log.Println(query)
	log.Println("----------------------------------------------------------")

	_, err := db.Exec(query, values...)
	return err
}
