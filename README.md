# Go Web Crawler

A concurrent web crawler written in Go that crawls pages from a single domain, extracts page content, and generates a structured JSON report containing information about each visited page.

## Features

* Concurrent crawling using Go goroutines
* Configurable concurrency limit
* Configurable maximum number of pages to crawl
* Restricts crawling to the same domain
* Prevents duplicate page visits
* Extracts page data from HTML
* Generates a `report.json` file with crawl results
* Thread-safe implementation using mutexes and wait groups

---

## Project Structure

```text
.
├── main.go              # Entry point
├── config.go            # Crawler configuration and crawling logic
├── report.json          # Generated crawl report
├── ...
```

Additional helper files (not shown) handle:

* HTML downloading
* URL normalization
* HTML link extraction
* Page content extraction
* JSON report generation

---

## Requirements

* Go 1.22+ (or newer)

---

## Installation

Clone the repository:

```bash
git clone https://github.com/yourusername/go-web-crawler.git
cd go-web-crawler
```

Install dependencies:

```bash
go mod tidy
```

---

## Usage

### Basic Usage

```bash
go run . https://example.com
```

This uses the default configuration:

* Maximum concurrent workers: **20**
* Maximum pages: **100**

---

### Specify Maximum Concurrency

```bash
go run . https://example.com 10
```

This limits crawling to **10 concurrent goroutines**.

---

### Specify Maximum Pages

```bash
go run . https://example.com 10 250
```

This configures:

* Maximum concurrency: **10**
* Maximum pages: **250**

---

## Command-Line Arguments

```text
go run . <website> [maxConcurrency] [maxPages]
```

| Argument       | Required | Default | Description                                |
| -------------- | -------- | ------- | ------------------------------------------ |
| website        | Yes      | —       | Website to crawl                           |
| maxConcurrency | No       | 20      | Maximum number of concurrent crawl workers |
| maxPages       | No       | 100     | Maximum number of pages to crawl           |

Example:

```bash
go run . https://golang.org 15 500
```

---

## How It Works

1. Parses the target URL.
2. Creates the crawler configuration.
3. Starts crawling from the root page.
4. Downloads HTML content.
5. Extracts page metadata and links.
6. Normalizes URLs to avoid duplicates.
7. Ignores links outside the target domain.
8. Spawns new goroutines for each discovered page.
9. Stops when the maximum page limit is reached.
10. Writes the collected results to `report.json`.

---

## Output

After the crawl completes, a JSON report is generated:

```text
report.json
```

The report contains the extracted information for every visited page.

Example structure:

```json
{
  "https://example.com": {
    "title": "...",
    "headings": [
      "..."
    ],
    "links": [
      "..."
    ]
  }
}
```

---

## Error Handling

The crawler validates:

* Missing website argument
* Invalid URLs
* Invalid numeric arguments
* HTML download failures
* URL parsing errors

Invalid input causes the program to terminate with a descriptive error message.

---

## Technologies Used

* Go
* Goroutines
* WaitGroups
* Mutexes
* Channels (Semaphore Pattern)
* HTML Parsing
* JSON Serialization

