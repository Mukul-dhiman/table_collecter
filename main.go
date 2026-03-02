package main

import (
	"fmt"
	"net/http"

	"Table_collecter/api"
	"Table_collecter/kafka"
)

func main() {
	brokers := []string{"localhost:9092"}
	topic := "tables"

	producer, err := kafka.NewProducer(brokers, topic)
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	api.InitKafka(producer)

	http.HandleFunc("/fetch-table", api.FetchTableHandler)

	fmt.Println("API running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
