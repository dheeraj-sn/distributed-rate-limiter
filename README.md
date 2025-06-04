# 🛠  Building a Distributed Rate Limiter in Go

## **🟢 Introduction**

Rate limiting is a critical part of modern infrastructure. It protects APIs, ensures fair usage, and prevents abuse — especially in multi-tenant, distributed systems.

I built **RLAAS (Rate Limiter as a Service)** to learn and demonstrate how a scalable, centralized rate limiter could work in Go. This post outlines the architecture, tradeoffs, and what I learned.

---

## **🚧 The Problem**

Imagine a login API or a payment service — without rate limiting, a single client could exhaust resources, DoS your app, or gain unfair advantage. Rate limiting solves this by enforcing a **policy of fairness and control**.

But doing this **across distributed stateless services** is tricky. In-memory counters don’t scale. Coordination is hard. Enter: **Redis-based centralized coordination**.

---

## **🧠 Design Goals**

- **Stateless API nodes** for horizontal scaling
- **Redis backend** for coordination and token state
- **Token Bucket** algorithm (burst-friendly)
- **In-memory fallback** when Redis is unavailable
- **Docker-first** development
- **gRPC-ready** and Kubernetes-ready foundation

---

## **🏗️ Architecture**

```jsx
Client
  |
  v
POST /check
  |
  v
[RLAAS Node]
  ├── RedisLimiter (primary)
  └── MemoryLimiter (fallback)
      |
      v
   [Redis Store]
```

- **RLAAS Node**: Accepts rate check requests.
- **RedisLimiter**: Uses atomic Redis operations or Lua to enforce limits.
- **MemoryLimiter**: Kicks in only when Redis is down (non-durable).

---

## **✅ MVP Features**

- /check endpoint:

```json
{
  "key": "user_123",
  "rate": 5,
  "interval": "1s"
}
```

- Redis-backed Token Bucket algorithm
- Graceful in-memory fallback
- Configurable per key
- Dockerized with Makefile
- Clean codebase for future extensibility

---

## **🧰 Tech Stack**

- **Go** 1.24
- **Redis** (central store)
- **HTTP API** (std/net + mux)
- **Docker**
- **GitHub Issues & Milestones** for planning

---

## **🧱 Code Highlights**

**Limiter Interface:**

```json
type Limiter interface {
    Allow(key string, rate int, interval time.Duration) bool
}
```

**Redis Limiter (Atomic):**

```json
// Uses INCR + EXPIRE or Lua scripts for atomic token logic
```

**Handler:**

```json
func (h *Handler) Check(w http.ResponseWriter, r *http.Request) {
    // parse request, call limiter.Allow, return JSON response
}
```

---

## **🔍 What I Learned**

- Designing **extensible interfaces** in Go
- Building for **failure modes** (Redis downtime)
- Writing **atomic** operations in Redis
- Structuring a Go service for **production readiness**
- Planning with **milestones, labels, GitHub CLI**

---

## **🚀 What’s Next**

- gRPC support
- Prometheus metrics
- CLI tool for testing
- Helm chart for Kubernetes
- Leaky bucket / sliding window plugins
- Redis clustering / Etcd backend