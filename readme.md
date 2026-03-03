# Table Collector

Table Collector is a modular ingestion system that fetches HTML tables from public URLs, normalizes them into consistent logical structures, streams them through Kafka, and persists them into MySQL.

It is designed to handle **real-world HTML**, including malformed tables, rowspans, layout noise, and inconsistent schemas (including Wikipedia pages).

---

## High-Level Architecture



Client

↓

HTTP API (FetchTableHandler)

↓

Fetcher (HTTP + gzip aware)

↓

Parser (table normalization via handlers)

↓

Kafka Producer

↓

Kafka

↓

Kafka Consumer

↓

DB Ingestion (dynamic schema + batch inserts)

↓

MySQL


Each layer has **one responsibility** and is intentionally decoupled from the others.

---

## Key Features

- HTTP fetcher with proper headers
- Table normalization independent of storage
- Multiple table handlers:
  - Simple tables
  - Rowspan-heavy tables
  - Key–value (vertical) tables
- Kafka-based producer–consumer decoupling
- Dynamic MySQL table creation
- Logical → physical type mapping in DB layer
- Batch inserts for performance
- Defensive parsing (bad HTML never crashes the server)
- End-to-end observability at producer and consumer boundaries

---

## Project Structure


Table_collecter/

├── api/ # HTTP handlers

├── fetcher/ # HTML fetching logic

├── parser/ # Table normalization + handlers

├── producer/ # Kafka producer

├── consumer/ # Kafka consumer


├── db/ # MySQL schema & insertion logic

├── monitor/ # System monitoring scripts

├── kafka_docker/ # Kafka docker setup

├── mysql_docker/ # MySQL docker setup

├── config/ # Runtime configs (ignored in git)

├── go.mod

├── go.sum

└── README.md


---



## Limitation

- No JS-rendered tables: Fails on SPAs, AJAX, or client-side content (basic HTTP only).

- No PDF support: Cannot extract from PDFs or embedded viewers (HTML DOM required).

- Limited colspan/rowspan: Complex merges cause misalignment or truncation.

- No auth-protected pages: Blocks login, OAuth, or session-based content.

- Append-only schemas: No versioning; ignores column changes or drops.