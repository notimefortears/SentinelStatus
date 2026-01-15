package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/notimefortears/sentinel/internal/monitor"
	"github.com/notimefortears/sentinel/internal/store"
	_ "github.com/lib/pq"
)

func main() {
	connStr := os.Getenv("DB_URL")
	var db *sql.DB
	var err error

	// Tracking failures in memory for alerting
	failureCount := make(map[string]int)

	for i := 0; i < 5; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
		}
		if err == nil {
			break
		}
		fmt.Printf("⚠️ Worker: DB not ready, retrying... (%d/5)\n", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal(err)
	}

	// Ensure tables exist
	db.Exec(store.Schema)

	// High-speed 5-second ticker
	ticker := time.NewTicker(5 * time.Second)
	fmt.Println("🛰️ SENTINEL_WORKER: HIGH_FREQUENCY_MODE_ENABLED")

	for range ticker.C {
		rows, err := db.Query("SELECT url FROM targets")
		if err != nil {
			fmt.Println("Error fetching targets:", err)
			continue
		}

		var urls []string
		for rows.Next() {
			var u string
			rows.Scan(&u)
			urls = append(urls, u)
		}
		rows.Close()

		if len(urls) == 0 {
			fmt.Println("😴 No targets in database. Waiting...")
			continue
		}

		fmt.Printf("\n[ %s ] ⚡️ FAST SCAN: %d targets\n", time.Now().Format("15:04:05"), len(urls))
		results := make(chan monitor.Result, len(urls))
		var wg sync.WaitGroup

		for _, url := range urls {
			wg.Add(1)
			go func(u string) {
				defer wg.Done()
				results <- monitor.CheckURL(u)
			}(url)
		}

		wg.Wait()
		close(results)

		for res := range results {
			// Log to Database
			_, err := db.Exec("INSERT INTO monitor_results (url, status_code, latency_ms) VALUES ($1, $2, $3)",
				res.URL, res.StatusCode, res.Latency.Milliseconds())
			
			if err != nil {
				fmt.Println("DB Insert Error:", err)
			}

			// --- HACKERMAN ALERTING LOGIC ---
			if res.StatusCode != 200 {
				failureCount[res.URL]++
				fmt.Printf("❌ FAIL [%d/3]: %s (Status: %d)\n", failureCount[res.URL], res.URL, res.StatusCode)
				
				if failureCount[res.URL] == 3 {
					fmt.Printf("🚨 ALERT: CRITICAL_FAILURE_DETECTED on %s\n", res.URL)
					fmt.Printf(">>> BROADCASTING ALERT FOR %s <<<\n", res.URL)
					// In a real prod environment, you'd trigger a Webhook/Slack call here
				}
			} else {
				// Reset failure count on successful check
				if failureCount[res.URL] > 0 {
					fmt.Printf("💚 RECOVERED: %s is back online\n", res.URL)
				}
				failureCount[res.URL] = 0
				fmt.Printf("✅ OK: %s (%dms)\n", res.URL, res.Latency.Milliseconds())
			}
		}
	}
}