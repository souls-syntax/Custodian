# CUSTODIAN
### *High-Performance Distributed Misinformation Detection System*

![Go](https://img.shields.io/badge/Backend-Go-blue)
![Python](https://img.shields.io/badge/ML-Python%20%2F%20BERT-yellow)
![Redis](https://img.shields.io/badge/Cache-Redis-red)
![Docker](https://img.shields.io/badge/Deploy-Docker%20Compose-2496ED)

---

Misinformation spreads in milliseconds, while verification takes minutes.

CUSTODIAN bridges this latency gap using a high-concurrency, tiered inference architecture. Built on a Go backend, it provides instant protection by resolving known rumors through a Redis cache in under a millisecond.

For new or unseen content, the system applies asymmetric confidence gating. A lightweight local DistilBERT model rapidly filters high-volume or obvious misinformation, while only ambiguous claims are escalated asynchronously to a cloud-based LLM for deeper reasoning. Human reviewers retain final authority over disputed cases by design.

By decoupling user response time from expensive AI processing, CUSTODIAN enables fast, scalable misinformation flagging while controlling cost and preserving accuracy.

## Key Features

* **Zero-Cost Infrastructure:** Runs entirely on local hardware behind a University firewall, using **Cloudflare Tunnels** to punch through CGNAT.
* **Tiered Inference:** * *Layer 1 (Fast):* Redis Cache (<1ms)
    * *Layer 2 (Medium):* Local DistilBERT Model (~50ms)
    * *Layer 3 (Slow/Deep):* Async LLM + Web Search Agent (~2s)
* **Anti-Fragile:** The system degrades gracefully. If the LLM is down, BERT takes over. If the Internet is down, the Cache takes over.
* **Eventual Consistency:** Users get an instant preliminary result (BERT), which is silently updated to a high-confidence result (LLM) in the background via the "Invisible Intelligence" pattern.

## System Architecture

1.  **Ingestion:** Query hits the **Go Orchestrator**.
2.  **Cache Layer:** Checks Redis. If hit -> **Return Instantly**.
3.  **Scatter-Gather:** If miss -> Query dispatched to **Postgres** (Audit) and **BERT Service** (Python) simultaneously.
4.  **Fast Path:** BERT analyzes the text. If confidence > threshold, result is returned and cached.
5.  **Async Escalation:** If BERT is uncertain, the job is pushed to a **Redis Queue**.
6.  **Deep Blue Worker:** A background worker picks up the job:
    * Optimizes the search term.
    * Queries live web sources.
    * Synthesizes a final verdict using Gemini LLM.
7.  **Write-Through:** The new, high-confidence verdict overwrites the cache and DB, ensuring the next user gets the "Verdict" instantly.

* **Backend:** Go (Golang) - *Chosen for high concurrency and low latency.*
* **ML/AI:** Python, DistilBERT, Gemini API, DuckDuckGo Search.
* **Data:** Redis (Cache), PostgreSQL (Persistent Storage).
* **Infrastructure:** Docker Compose, Cloudflare Zero Trust.


## How to run

### Prerequisites
* Docker & Docker Compose installed.
* A `.env` file with `GOOGLE_API_KEY` and `POSTGRES_PASSWORD`.

### Deploy
```bash
# 1. Clone the repository
git clone https://github.com/souls-syntax/Custodian

# 2. Start the build
docker-compose up -d --build

# 3. Verify status
docker ps
```
## Limitations

As with all statistical models, CUSTODIAN does not produce absolute truth.  
Its outputs should be treated as probabilistic assessments designed to aid verification, not replace human judgment.
