package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Target represents the input for adding/deleting URLs
type Target struct {
	URL string `json:"url"`
}

// Stat represents the data sent to the dashboard
type Stat struct {
	URL        string `json:"url"`
	StatusCode int    `json:"status_code"`
	Latency    int    `json:"latency_ms"`
	CheckedAt  string `json:"checked_at"`
}

var db *sql.DB

func main() {
	var err error
	connStr := os.Getenv("DB_URL")

	// Robust connection logic for Docker startup
	for i := 0; i < 5; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
		}
		if err == nil {
			break
		}
		fmt.Printf("⚠️ API: DB not ready, retrying... (%d/5)\n", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("Could not connect to database:", err)
	}

	// 1. SERVE DASHBOARD: Serves everything in the /web folder
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	// 2. TARGETS ENDPOINT: Manage the list of URLs
	http.HandleFunc("/targets", handleTargets)

	// 3. STATS ENDPOINT: Provide data to the dashboard
	http.HandleFunc("/api/stats", handleStats)

	fmt.Println("🚀 Sentinel Manager & Dashboard running on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func handleTargets(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var t Target
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil || t.URL == "" {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}
		_, err := db.Exec("INSERT INTO targets (url) VALUES ($1) ON CONFLICT DO NOTHING", t.URL)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Target added: %s", t.URL)

	case "GET":
		rows, _ := db.Query("SELECT url FROM targets")
		var urls []string
		for rows.Next() {
			var u string
			rows.Scan(&u)
			urls = append(urls, u)
		}
		json.NewEncoder(w).Encode(urls)

	case "DELETE":
		var t Target
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil || t.URL == "" {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}
		_, err := db.Exec("DELETE FROM targets WHERE url = $1", t.URL)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		fmt.Fprintf(w, "Target deleted: %s", t.URL)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleStats(w http.ResponseWriter, r *http.Request) {
    stats := []Stat{} // Initialize as empty slice so we return [] not null

    query := `
        SELECT url, status_code, latency_ms, checked_at FROM (
            SELECT url, status_code, latency_ms, checked_at,
            ROW_NUMBER() OVER (PARTITION BY url ORDER BY checked_at DESC) as rn
            FROM monitor_results
        ) t
        WHERE rn <= 20
        ORDER BY url ASC, checked_at DESC`

    rows, err := db.Query(query)
    if err != nil {
        // HACKERMAN LOGIC: If table doesn't exist, just return an empty list 
        // instead of a 500 error.
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(stats)
        return
    }
    defer rows.Close()

    for rows.Next() {
        var s Stat
        var t time.Time
        if err := rows.Scan(&s.URL, &s.StatusCode, &s.Latency, &t); err != nil {
            continue
        }
        s.CheckedAt = t.Format("15:04:05")
        stats = append(stats, s)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(stats)
}