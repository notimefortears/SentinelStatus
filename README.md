# 🛰️ SENTINEL // Distributed Monitoring System

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)


A high-performance, distributed site reliability monitor. Sentinel provides real-time latency tracking, consecutive failure alerting, and a time-series telemetry dashboard—all wrapped in a professional, dark-mode terminal interface.



---

## ⚡ Key Features

* **Distributed Architecture:** Decoupled Manager API and Worker nodes.
* **High-Frequency Polling:** Sub-5-second latency checks via Go routines.
* **Time-Series Analytics:** Historical performance trends visualized via Chart.js.
* **Self-Healing Logic:** Consecutive failure thresholds (3/3) to prevent false positives.
* **Dynamic Badges:** Auto-generated SVG status badges for external embedding.
* **Containerized:** One-command deployment using Docker Compose.

---

## 🏗️ System Architecture



1.  **Manager API (Go):** The control plane for target CRUD operations and UI hosting.
2.  **Sentinel Worker (Go):** Concurrent monitoring engine using `sync.WaitGroup` and channels.
3.  **PostgreSQL Store:** Relational database using Window Functions (`PARTITION BY`) for trend analysis.
4.  **Telemetry UI:** A lightweight Vanilla JS frontend for real-time data visualization.

---

## 🚀 Deployment

### Prerequisites
- Docker & Docker Compose
- Make (optional, but recommended)

### Quick Start
```bash
# 1. Spin up the infrastructure
make start

# 2. Inject your first monitoring target
curl -X POST http://localhost:3000/targets \
     -H "Content-Type: application/json" \
     -d '{"url": "[https://google.com](https://google.com)"}'

# 3. Access the Terminal Interface
# URL: http://localhost:3000

##🛠️ Operational Commands

###The system is fully automated via the Makefile. Run these commands from the root directory:

Action	Command	Command Description
Initialize	make start	Build images & spin up api, worker, db, and adminer.
Shutdown	make stop	Stop all services without deleting data.
Live Logs	make logs	Stream the Sentinel Worker's heartbeat and check results.
Reset	make clean	Destructive: Wipe database volumes and clear Docker cache.
Seed	make test-target	Automatically inject google.com into the monitor loop.
##📡 API Documentation

The Sentinel API serves as the orchestration layer. Use the following endpoints to interact with the system programmatically:
CORE Endpoints

    Add Target POST /targets — Body: {"url": "https://example.com"}

    Delete Target DELETE /targets — Body: {"url": "https://example.com"}

    List All GET /targets — Returns array of strings.

TELEMETRY Endpoints

    Historical Stats GET /api/stats

    Returns the last 20 checks per URL including latency_ms and status_code. Optimized via Postgres Window Functions.

    Dynamic Badge GET /api/badge?url={URL}

    Returns a live-updating SVG.

    Usage: <img src="http://localhost:3000/api/badge?url=https://google.com" />

###📸 Dashboard Preview

The frontend is a Single Page Application (SPA) designed for NOC (Network Operations Center) displays.
Feature	Specification
Theme	Monochrome / Neon Green / Dark Mode
Updates	Real-time (2000ms polling)
Charts	Bézier Curve Latency Trends (Chart.js)
Mobile	Fully Responsive CSS Grid
###⚙️ Project Configuration
.env Structure

The system expects the following environment variables (handled automatically by docker-compose):
Bash

DB_URL=postgresql://admin:password@db:5432/sentinel?sslmode=disable
POSTGRES_USER=admin
POSTGRES_PASSWORD=password

