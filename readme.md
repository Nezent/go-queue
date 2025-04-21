# 🚀 GoQueue – Scalable Job Queue & Task Processing System in Golang

GoQueue is a production-ready, event-driven Job Queue system built with **Golang** and **PostgreSQL**. It’s designed to teach you how to build clean, secure, and scalable backends — with background jobs, scheduling, retries, and real-time updates.

---

## 🧠 Project Objective

> Build a secure, multi-user job queue system where authenticated users can enqueue tasks, track their job statuses, and receive real-time updates when jobs complete — all in pure Go + PostgreSQL.

---

## 📦 Project Structure (Hexagonal Architecture)

![Hexagonal Architecture](./hexagonal_architecture.webp)

```
go-queue/
├── cmd/
│   ├── api/                # HTTP server (REST API)
│   └── worker/             # Background worker runner
├── internal/
│   ├── domain/             # Job/User models, interfaces
│   ├── service/            # Job logic (enqueue, execute, retry)
│   └── repository/         # PostgreSQL implementations
├── migrations/             # SQL files for DB setup
├── common/                 # Shared utilities (JWT, auth)
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── README.md
```

---

## 🎯 Application Features Overview

This job queue system is designed with **security**, **performance**, and **developer-friendliness** in mind. Below is a summary of what it offers:

### 🔐 Secure Multi-User System
- Users can **sign up** and **log in**
- Passwords are hashed using `bcrypt`
- JWT-based authentication with access tokens
- Middleware to protect routes and identify the current user
- Users can only manage **their own jobs**

### 📦 Job Management
- Submit new jobs with type, payload, and priority
- Job queue system manages status: `pending`, `processing`, `completed`, `failed`
- Users can **retry**, **cancel**, or **delete** their own jobs
- Each job tracks number of attempts, timestamps, and status updates

### ⚙️ Background Processing
- Lightweight **Goroutines** handle job execution asynchronously
- Retry failed jobs automatically based on logic
- Jobs can be scheduled for future execution using `run_at` field

### 📡 Real-Time Updates (WebSocket)
- WebSocket connection to push live job updates to the user dashboard
- Instant status refresh for running/completed/failed jobs

### 🛡️ Access Control & Isolation
- Users are sandboxed; no cross-access to other users' jobs
- Auth middleware ensures only valid users can access protected routes

### 📊 Developer & DevOps Friendly
- Modular code structure using **Go + PostgreSQL**
- RESTful API built with **Fiber** or **Chi**
- Uses **pgx** or **SQLx** for database layer
- Dockerized for easy local dev and deployment
- Can be deployed to **GCP**, **Fly.io**, etc.

---

## 🚦 Roadmap & Learning Milestones

### ✅ Phase 1: Foundation – REST API + PostgreSQL + Auth

- [x] Set up project structure with Hexagonal Architecture
- [ ] Create `users` and `jobs` tables
- [ ] PostgreSQL + pgx
- [ ] Add user authentication:
  - [ ] Signup (`POST /signup`)
  - [ ] Login (`POST /login`)
  - [ ] JWT token generation & middleware
- [ ] REST API for:
  - [ ] Submit job (auth required)
  - [ ] Get job status (auth required)
- [ ] Docker + Compose setup

---

### ✅ Phase 2: Worker System – Background Task Execution

- [ ] Goroutine-based worker
- [ ] Poll for pending jobs
- [ ] Execute job logic (mocked at first)
- [ ] Retry with backoff
- [ ] Log output & mark as completed/failed

---

### ✅ Phase 3: Scheduling + Priority Queue

- [ ] `run_at` timestamp support
- [ ] `priority` field (`high`, `medium`, `low`)
- [ ] Queue sorted by run_at + priority

---

### ✅ Phase 4: Dashboard + Monitoring

- [ ] Web UI (or CLI) to view jobs by user
- [ ] Retry failed jobs manually
- [ ] Add `/metrics` endpoint for Prometheus
- [ ] Optional: Grafana setup

---

### ✅ Phase 5: WebSocket Notifications (Advanced)

- [ ] Notify logged-in users in real-time when their job completes
- [ ] WebSocket connection with JWT auth
- [ ] Frontend toast/alert when status updates

---

### ✅ Phase 6: Advanced Concepts

- [ ] Dead Letter Queue (failed jobs after N retries)
- [ ] Graceful shutdown with context & signals
- [ ] Redis cache or pub/sub (optional)
- [ ] Deploy to GCP or any cloud
- [ ] CI/CD (GitHub Actions)

---

## 🔐 Database Schema (Users + Jobs)

![Database Schema](./migrations/database_schema.png)
---

## 🌐 API Endpoints

### 🔑 Auth Routes
- `POST /signup` – Register new user
- `POST /login` – Authenticate and get token

### 📦 Job Routes (require JWT)
- `POST /jobs` – Submit a new job (auth required)
- `GET /jobs/:id` – View job (only if you own it)
- `GET /jobs` – List your jobs (by status, priority)
- `POST /jobs/:id/retry` – Retry failed job

### 📡 Real-Time
- `WS /ws/jobs` – Connect with JWT, get updates

---

## 🚀 Tech Stack

- 🧠 **Language:** Go (Golang)  
- 🗄️ **Database:** PostgreSQL  
- 🔐 **Authentication:** JWT + bcrypt  
- 🌐 **API Framework:** Chi  
- 🛢️ **DB Layer:** pgx  
- 🧵 **Background Tasks:** Goroutines  
- 🐳 **DevOps & Containerization:** Docker  
- ⚡ **Realtime Communication:** WebSocket  
- ☁️ **Deployment:** Github

---

## 🧠 Learnings From This Project

- ✅ Clean, modular Go architecture
- ✅ Job queueing with retries and scheduling
- ✅ Auth-secured APIs with JWT
- ✅ PostgreSQL for relational & task data
- ✅ Real-time systems with WebSocket
- ✅ Dockerized microservices
- ✅ Observability & cloud-readiness

---

## 📚 Resources

- [JWT Authentication in Go](https://dev.to/macisamuele/jwt-authentication-in-go-1j7h)
- [Hexagonal Architecture in Go](https://medium.com/@matryer/structuring-go-applications-clean-architecture-ef7d7c6fcd26)
- [Go Channels and Workers](https://gobyexample.com/worker-pools)
- [Job Queues in Postgres](https://www.crunchydata.com/blog/building-a-job-queue-with-postgresql)

---

## 🤝 Contributors

Made with ❤️ by Sirajum Munir (Nezent)


