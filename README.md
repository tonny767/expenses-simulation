# Expense Management System

A backend-focused Expense Management System built as a take-home challenge.

---

## 1. Setup & Run Instructions

### Prerequisites

* Go 1.21+
* Docker & Docker Compose
* PostgreSQL 18
* Nuxt.js 4

### Environment Variables

Create a `.env` file at the project root or copy from .env.example
```bash
cp .env.example .env
```
modify it according to your needs e.g. (DATABASE_URL/BACKEND_URL)

### Run with Docker Compose
Download Docker Desktop (simplest) or make sure docker works by doing:
```bash
docker info
```
then you can do:
```bash
docker compose up --build
```


you should see all three services started (expense-backend, expense-frontend and expense-postgres)

Notes: 
- Make sure backend url is `http://backend:8080` because of the naming of the docker
- Change the database values in both `DATABASE_URL` and parts of `DB_`, look at the log if there is error.

Services started:

* API server: `http://backend:8080` for Docker or `http://localhost:8080` for manual setup
* PostgreSQL
* 
* Background payment worker

---
### Run with Manual Setup

Make sure you have the required packages ready.

#### GoLang setup

Go to backend

`cd backend`

Run this to create the `go.sum`

`go mod tidy`

Install the GoLang 

`go install`





---

## 2. API Documentation

Swagger UI is available at:

```
http://localhost:8080/swagger/index.html
```

### Currency Rules

* All amounts are **IDR only**
* Example:

```json
{
  "amount_idr": 150000,
  "description": "Client meeting lunch",
  "receipt_url": "/receipts/lunch.png"
}
```

---

## 3. Expense Visibility & Roles

### User

* Can only see **their own expenses**
* Can view:

  * Expense details
  * Expense status
  * Approval status

### Manager

* Can see **all expenses**
* Can view:

  * Expense owner
  * Approval records
* Responsible for approving or rejecting expenses

---

## 4. Business Rules Implementation

### Approval Threshold

* Expenses below **IDR 1,000,000** are auto-approved
* Expenses above the threshold require manager approval

### Status Flow

```
PENDING
  ↓ (manager approves)
APPROVED
  ↓ (payment worker)
PROCESSING
  ↓ (success)
COMPLETED
```

Rejection flow:

```
PENDING → REJECTED
```

### Key Rules

* Approval must complete before payment starts
* Payment is asynchronous
* Expense remains `APPROVED` until payment succeeds

---

## 5. Architecture Decisions

### Clean Architecture

Development order:

1. Domains
2. Business rules
3. Actions (use cases)
4. Migrations
5. Repositories
6. API layer

This ensures business logic is independent of HTTP and frameworks.

### UUID for External Access

* Internal database uses numeric IDs
* External APIs expose `external_id` (UUID)

**Reason:** Prevent ID enumeration and improve security.

### Background Worker

* Payment processing handled asynchronously
* Uses idempotency safeguards

---

## 6. Assumptions

* Currency is IDR only
* Single approval level
* Expense data is sufficiently displayed in table view
* Backend operates fully in UTC
* JWT stored in HTTP-only cookies

---

## 7. Limitations

* Status transition rules not fully enforced
* No approval audit log
* No retry/backoff for failed payments
* Worker restart is manual via API
* No rate limiting implemented

---

## 8. Trade-offs

### Approval Before Payment

Ensures financial correctness but increases state complexity.

### No Expense Detail Page

Reduces UI complexity; table already contains all relevant information.

### Worker Failure Handling

If worker fails, expense remains `APPROVED` and requires manual retry.

---

## 9. Improvements With More Time

### Architecture

* Explicit state machine for expense transitions
* Domain-level transition guards

### Features

* Multi-level approvals
* Approval audit logs
* Retry strategy with backoff
* Dead-letter queue

### DevOps

* Worker health checks
* Auto-restart workers
* Metrics and structured logging

---

## 10. Final Note

This system prioritizes **financial correctness, security, and extensibility** over feature completeness. Certain UX and automation improvements were intentionally deferred to ensure a solid and maintainable core architecture.


initially have auto-approved status, but removed because not needed.