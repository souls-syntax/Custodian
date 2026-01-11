# CUSTODIAN
### Missinformation detection system

---

Misinformation spreads in milliseconds, while verification takes minutes.

CUSTODIAN bridges this latency gap using a high-concurrency, tiered inference architecture.
Built on a Go backend, it provides instant protection by resolving known rumors through a Redis cache in under a millisecond.
For new content, the system applies asymmetric confidence gating: a lightweight local TinyBERT model rapidly filters high-volume
or obvious misinformation, while only ambiguous claims are escalated asynchronously to a cloud-based LLM for deeper reasoning. Human 
reviewers retain final authority over disputed cases. By decoupling user response time from expensive AI processing, CUSTODIAN enables fast, 
scalable misinformation flagging while controlling cost and preserving accuracy.
