# dd-go-api

![Go Version](https://img.shields.io/badge/go-1.20+-blue)
![Dockerized](https://img.shields.io/badge/docker-ready-blue)
![Issues](https://img.shields.io/github/issues/RecursionExcursion/dd-go-api)
![Last Commit](https://img.shields.io/github/last-commit/RecursionExcursion/dd-go-api)

An API written in Go for [Dune Digital](https://dunedigital.io), split into two main domains: **WSD (Workspace Deployer)** and **BetBot**.

---

## Domains

### WSD - Workspace Deployer

WSD is a web API that accepts scripts and returns binaries executable on the client’s system to set up their workspace. This includes launching URLs, applications, and more to automate environment bootstrapping.

### BetBot

BetBot is a data aggregation API that scrapes thousands of ESPN API endpoints to compile real-time sports stats for betting analytics and predictions.

---

## Tech Stack

- **Language**: Go (Golang)
- **Containerized**: Docker

---

## Getting Started

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
BB_JWT_SECRET=<your_bb_jwt_secret>

WSD_API_KEY=<your_wsd_api_key>

SELF_URL=<hosted_url>

CFB_API_KEY=<ctfr_api_key>
DB_NAME_CFBR=<cfbr_db_name>

DB_NAME_PICKLE=<your_pickle_db_name>
PICKLE_USERNAME=<app_username>
PICKLE_PASSWORD=<app_password>
PICKLE_SECRET=<app_secret>
```

### 4. Run

### Option A- Run Locally

```bash
go run .
```

### Option B- Run with Docker

```bash
docker build -t dd-go-api .
docker run -p 8080:8080 --env-file .env dd-go-api
```

## Endpoints Overview

- /wsd/\* – All routes related to the Workspace Deployer

- /betbot/\* – All routes for BetBot stats and data access

- /cfbr/\* - All routes for college football ranker

- /pickle\* - All routes for The Pickle

- For a full list of available routes, see the docs folder or inspect routes.go.

## Notes

- APIs require domain-specific API keys for access.

- Some routes uses JWT for authentication.
