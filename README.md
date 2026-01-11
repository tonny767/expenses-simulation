# Expense Management System

A backend-focused Expense Management System built as a take-home challenge.

---

## Technology Stack

**Backend:**
- Go 1.23 + Gin Framework
- GORM (PostgreSQL ORM)
- JWT Authentication
- Background Workers (Goroutines)

**Frontend:**
- Nuxt.js 4 + Vue 3 (Composition API)
- Shadcn-vue + TailwindCSS

**Database:**
- PostgreSQL 18

**DevOps:**
- Docker & Docker Compose
- Goose Migrations
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

copy .env 
```
copy .env.example .env
```

then you can do:
```bash
docker compose up --build
```

For the data seeding, make sure all three services works first, then:

```bash
docker compose exec backend sh
```

Then run this command to seed the data, make sure your DATABASE_URL in .env is correct first
```
goose -dir ./migrations postgres "$DATABASE_URL" up
```

you should see all three services started (expense-backend, expense-frontend and expense-postgres)

Notes: 
- Make sure backend url is `http://backend:8080` because of the naming of the docker
- Change the database values in both `DATABASE_URL` and parts of `DB_`, look at the log if there is error.

Services started:

* Backend API server: `http://backend:8080` for Docker or `http://localhost:8080` for manual setup
* Database PostgreSQL
* Frontend Nuxt/Vue.js
* Background payment worker

---
### Run with Manual Setup

Make sure you have the required packages ready.

### Postgresql setup
If you have psql ready or running from Docker can skip this part, in MacOS run this:
```bash
brew install postgres@18
```

make sure psql runs, can also determine your own User by e.g. `psql -U postgres -H postgres`
```bash
psql
```
inside run commands to add DB
```bash
CREATE DATABASE expenses;
```

then `\q` and you should already have your database ready, don't forget to modify your `.env` according to the database you initialized above.

---

### GoLang setup

Go to backend

```bash
cd backend
```

Run this to create the `go.sum`
```bash
go mod tidy
```

Install the GoLang 
```bash
go mod download
```

Install `goose` for migrations 
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

run the database seeding, make sure you have db ready and `.env` written. includes user and sample expenses 
```bash
goose -dir ./migrations postgres "$DATABASE_URL" up
```

now you can run the `main.go`
```bash
go run main.go
```
---

### Nuxt Setup

```bash
cd frontend
```
if you can `npm` can just do
```bash
npm install
```
then you can run:
```bash
npx nuxi prepare

```
to run the frontend, simply do: 
```bash
npx nuxi dev
```

---

## 2. API Documentation

To see all the available GET/POST/PUT functions
Swagger UI is available at:

```
http://localhost:3000/v1/api/swagger/index.html
```

### Currency Rules

* All amounts are **IDR only**
* Example:

```json
{
  "amount_idr": 150000, // Rp 150.000
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
* Can create:
  * Expense

### Manager

* Can see **all expenses**
* Can view:
  * Expense audit logs
  * Pending Approvals record
  * All user's expenses record
* Can approve/reject expenses


---
## 4. Business Rules Implementation

### Approval Threshold

* Expenses below **Rp 1.000.000** are auto-approved
* Expenses above the threshold require manager approval
* Limit is Minimum of **Rp. 10.000** and Maximum of **Rp. 50.000.000** 

### Status Flow

Auto-approved flow:
```
PENDING
  ↓ (auto-approved)
APPROVED
  ↓ (payment worker)
COMPLETED
```

Approval required flow:
```
PENDING
  ↓ (manager approves)
APPROVED
  ↓ (payment worker)
COMPLETED
```

Rejection flow:

```
PENDING → REJECTED
```

### Key Rules

* Approval must complete before payment starts
* Payment is asynchronous (refresh page to see changes on status after approving or creating ~4-5 seconds)
* Expense remains `APPROVED` until payment succeeds

---

## 5. Architecture Decisions

### Clean Architecture

Development order using GIN+GORM as framework:

1. Domains
2. Business rules
3. Actions (use cases)
4. Migrations
5. Controllers
6. API layer

This ensures business logic is independent of HTTP and frameworks.

### UUID for External Access

* Internal database uses numeric IDs
* External APIs expose `external_id` (UUID)

**Reason:** Prevent ID enumeration and improve security.

### Background Worker

* Payment processing handled asynchronously
* Uses idempotency safeguards (currently its just mock so real one is not yet exist)

---

## 6. Assumptions

* Currency is IDR only
* Single approval level
* Expense data is sufficiently displayed in table view
* Backend operates fully in UTC, frontend will translate to `Asia/Jakarta`
* JWT stored in HTTP-only cookies
* Receipt URL is hard-coded for now
* Mock payment that prevents idempotency is yet to work
* Users are pre-seeded, no register required
* Status transition is enforced in `canTransition` rules.
* Right now approval does not get shown
* Removed auto-approved as status because its redundant
* Reverse Proxy is applied in `/v1/api` format
* Expense and Approval must have same statuses across all the processes (unless `COMPLETED` flag)

---

## 7. Limitations

* No approval audit log (in the approvals part)
* No rate limiting implemented
* No production-ready environment 
* No manual entry API to retry the payment process

---

## 8. Trade-offs

### Approval Before Payment

Ensures financial correctness but increases state complexity.

### Worker Failure Handling

If worker fails, expense and approval remains `APPROVED` and requires manual modification.

### Pre-seeded accounts

No registration section

### No email notifications

For this feature, I need to enable approval flow and set manager first before creating approval, not possible with current approach unless I hard-coded the approver.

---

## 9. Improvements With More Time

### Architecture

* Explicit state machine for expense transitions
* Domain-level transition guards
* Repository implementation

### Features

* Multi-level approvals
* More audit logs showing e.g login
* Processing flag for expenses
* Dead-letter queue
* Login and Logout approach more beautifully
* Toaster instead of alert
* Idempotency table to limit user submitting
* Using UUID to do explicit API calls like accessing details page or approve/rejecting expenses and make sure no ID is exposed in the website (can hold security risk)
* User delete `PENDING` expenses

### DevOps

* Worker health checks
* Auto-restart workers
* Swagger
* Docker-ready environment
* Metrics and structured logging

---
## 10. Testing
In Docker `docker compose exec backend sh` in manual setup `cd backend`

Go to tests folder
```
cd tests
```

Run the test
```
go test
```

PASSED will be the final result
The test covered all the actions 
* Auto-approved Expense
* Pending Approval Expense
* Invalid Amount Expense
* Approve Expense
* Reject Expense
* Running mock payment process in auto-approved and approved expenses

---
## 11. Final Note

This system prioritizes **financial correctness, security, and extensibility** over feature completeness. Certain UX and automation improvements were intentionally deferred to ensure a solid and maintainable core architecture.

Login Details 
```bash
// Manager
alice@manager.com
password123

// User
bob@user.com
password123

dave@user.com
password123

eve@user.com
password123
```
