package api

import (
	"log"
	"net/http"

	"Table_collecter/fetcher"
	"Table_collecter/kafka"
	"Table_collecter/parser"
)

var producer *kafka.Producer

func InitKafka(p *kafka.Producer) {
	producer = p
}

func FetchTableHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "url parameter required", http.StatusBadRequest)
		return
	}

	log.Println("FETCH START:", url)

	// 1️⃣ Fetch
	doc, err := fetcher.FetchDocument(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("FETCH OK")

	// 2️⃣ Parse
	result, err := parser.ParseTables(doc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("PARSE OK: %d tables found\n", len(result.Tables))

	// 🔍 OPTIONAL: print full parsed tables
	// debugPrintResult(result)

	// 3️⃣ Enqueue
	if err := producer.Enqueue(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("ENQUEUE OK")

	w.Write([]byte("Tables enqueued to Kafka"))
}
