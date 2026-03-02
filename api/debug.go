package api

import (
	"Table_collecter/parser"
	"log"
)

func debugPrintResult(result *parser.Result) {
	log.Println("========== PARSED TABLES ==========")
	log.Printf("Total tables: %d\n", len(result.Tables))

	for _, t := range result.Tables {
		log.Printf("TABLE INDEX: %d\n", t.Index)

		log.Println("HEADERS:")
		for i, c := range t.Columns {
			log.Printf("  [%d] %s (%s)\n", i, c.Name, c.Type)
		}

		log.Println("ROWS:")
		for rIdx, row := range t.Rows {
			log.Printf("  ROW %d: %v\n", rIdx, row)
		}

		log.Println("-----------------------------------")
	}

	log.Println("===================================")
}
