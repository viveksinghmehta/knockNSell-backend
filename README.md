# 🏎️ knockNSell-backend

A high-performance, scalable backend for **KnockNSell** — built using Go, Gin, PostgreSQL, and AWS Lambda.

This repository contains the core backend logic, API routes, and database layer for the KnockNSell application.

---

## 🚀 Tech Stack

- **Go** (Golang)
- **Gin** web framework
- **PostgreSQL**
- **AWS Lambda**
- **SQLC** for type-safe database access
- **GitHub Actions** for CI/CD
- **Logrus** for structured logging

---

## 📂 Project Structure

<pre>
knockNSell-Backend
├── go.mod
├── go.sum
├── Makefile
├── README.md
├── sqlc.yml
├── .github/
│   └── workflows/
│       └── go.yml
├── cmd/
│   └── main.go                # Main entry point for Lambda & Debug mode
├── db/
│   ├── gen/                   # SQLC generated code
│   ├── queries/               # SQL query definitions
│   ├── schema/                # Database schema migrations
├── helpers/
│   ├── logflareLogger.go      # Logflare / Discord / Log setup
│   ├── SetupRouterAndLogger.go
│   ├── slackLogger.go
│   ├── sqlHelper.go
│   └── tokenHelper.go
├── routes/
│   ├── loginSignup.go
│   ├── ping.go
│   ├── router.go              # Central router registration
│   ├── sendOTP.go
│   ├── server.go              # Server struct for handler dependencies
│   ├── updateProfile.go
│   └── verifyOTP.go
</pre>

---

## 🛠️ Setup & Run Locally

### Prerequisites

- Go 1.24+
- PostgreSQL
- [sqlc](https://sqlc.dev/) installed
- AWS CLI configured (if deploying)

### 1️⃣ Clone the repo

```bash
git clone https://github.com/viveksinghmehta/knockNSell-backend.git
cd knockNSell-backend
