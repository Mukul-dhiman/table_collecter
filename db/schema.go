package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

func mysqlType(logical string) string {
	switch strings.ToUpper(logical) {
	case "STRING":
		return "VARCHAR(255)"
	case "VARCHAR":
		return "VARCHAR(255)"
	case "INT":
		return "INT"
	case "FLOAT":
		return "DOUBLE"
	case "BOOL":
		return "BOOLEAN"
	default:
		return "VARCHAR(255)"
	}
}

func EnsureTable(db *sql.DB, tableName string, columns []Column) error {
	var defs []string
	seen := map[string]bool{}

	for i, c := range columns {
		log.Printf("RAW COLUMN [%d]: '%s'\n", i, c.Name)

		colName := SanitizeColumnName(c.Name, i)
		log.Printf("SANITIZED COLUMN [%d]: '%s'\n", i, colName)

		if colName == "" {
			continue
		}

		if seen[colName] {
			colName = fmt.Sprintf("%s_%d", colName, i)
		}
		seen[colName] = true

		sqlType := mysqlType(c.Type)

		defs = append(defs, fmt.Sprintf("`%s` %s", colName, sqlType))
	}

	if len(defs) == 0 {
		return fmt.Errorf("no valid columns after sanitization")
	}

	query := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS `%s` (id INT AUTO_INCREMENT PRIMARY KEY, %s)",
		tableName,
		strings.Join(defs, ", "),
	)

	log.Println("FINAL CREATE TABLE QUERY:")
	log.Println(query)
	log.Println("===============================================================")

	_, err := db.Exec(query)
	return err
}
