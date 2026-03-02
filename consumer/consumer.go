package main

import (
	"encoding/json"
	"fmt"
	"log"

	"Table_collecter/db"
	"Table_collecter/kafka"
)

type Payload struct {
	Tables []db.Table `json:"tables"`
}

func main() {
	// Load DB config
	cfg, err := db.LoadConfig("../config/db.yaml")
	if err != nil {
		panic(err)
	}

	mysql, err := db.NewMySQL(cfg)
	if err != nil {
		panic(err)
	}
	defer mysql.Close()

	consumer, err := kafka.NewConsumer(
		[]string{"127.0.0.1:9092"},
		"tables",
	)
	if err != nil {
		panic(err)
	}

	log.Println("Kafka → MySQL consumer started")

	consumer.Consume(func(msg []byte) {
		var payload Payload
		if err := json.Unmarshal(msg, &payload); err != nil {
			log.Println("Invalid payload:", err)
			return
		}

		for _, t := range payload.Tables {
			tableName := "scraped_table_" + fmt.Sprint(t.Index)

			if err := db.EnsureTable(mysql, tableName, t.Columns); err != nil {
				log.Println("Table create failed:", err)
				continue
			}

			if err := db.InsertRows(mysql, tableName, t.Columns, t.Rows); err != nil {
				log.Println("Insert failed:", err)
				continue
			}

			log.Println("Inserted rows into", tableName)
		}
	})
}
