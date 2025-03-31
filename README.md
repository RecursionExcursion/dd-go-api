# dd-go-api

An API written in Go for [Dune Digital](https://dunedigital.io), split into two main domains: **WSD (Workspace Deployer)** and **BetBot**.

---

## 🧩 Domains

### 🔧 WSD - Workspace Deployer
WSD is a web API that accepts scripts and returns binaries executable on the client’s system to set up their workspace. This includes launching URLs, applications, and more to automate environment bootstrapping.

### 🤖 BetBot
BetBot is a data aggregation API that scrapes thousands of ESPN API endpoints to compile real-time sports stats for betting analytics and predictions.

---

## 🚀 Tech Stack

- **Language**: Go (Golang)
- **Containerized**: Docker

---

## ⚙️ Getting Started

### 1. Clone the repo

```bash
git clone https://github.com/your-org/dd-go-api.git
cd dd-go-api
```

### 2. Install Go modules

```bash
go mod tidy
```

### 3. Configure .env

```bash
PORT=<int>
ATLAS_URI=<your_mongo_connection_uri>
DB_NAME_BB=<your_betbot_db_name>
BB_API_KEY=<your_bb_api_key>
WSD_API_KEY=<your_wsd_api_key>
LOG_LEVEL=<int>
BB_JWT_SECRET=<your_jwt_secret>
```

### 4. Run

## Option A- Run Locally

```bash
go run .
```

## Option B- Run with Docker


```bash
docker build -t dd-go-api .
docker run -p 8080:8080 --env-file .env dd-go-api
```

## 📦 Endpoints Overview

/wsd/* – All routes related to the Workspace Deployer

/betbot/* – All routes for BetBot stats and data access

For a full list of available routes, see the docs folder or inspect routes.go.

## 🧠 Notes

Logging verbosity is controlled via LOG_LEVEL.

APIs require domain-specific API keys for access (WSD_API_KEY and BB_API_KEY).

BetBot uses JWT for authentication – keep BB_JWT_SECRET secure.