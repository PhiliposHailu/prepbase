# PrepBase API

PrepBase is a high-performance, RESTful backend engine designed for developers to share, discuss, and study technical interview questions.

It is built strictly using **Go (Golang)**, the **Gin Web Framework**, and **MongoDB**, adhering entirely to **Clean Architecture (Onion Architecture)** principles for maximum scalability and testability.

---

## 🏗️ Architectural Highlights

- **Strict Clean Architecture:** The application is divided into `Domain`, `Usecase`, `Repository`, `Delivery`, and `Infrastructure` layers. The core business logic has zero dependencies on web frameworks or databases.
- **Dependency Injection:** All external services (MongoDB, JWTs, Password Hashing, AI) are abstracted behind interfaces and injected at runtime, allowing for instant mocking during testing.
- **Concurrency (Goroutines & Channels):**
  - AI Hint generation (`POST /questions/:id/hint`) is processed concurrently.
  - A background Garbage Collector Goroutine sweeps the in-memory cache every 5 minutes to prevent memory leaks.
- **Asynchronous SMTP:** "Forgot Password" emails are dispatched via Mailtrap.io in a background Goroutine to ensure sub-millisecond API response times.
- **Thread-Safe Caching:** Custom implementation of a TTL-based caching service using `sync.RWMutex` to protect against race conditions.
- **Database Indexing:** Leverages MongoDB indexes to accelerate cursor-based pagination and search filters, while using Compound Unique Indexes to enforce strict data integrity (e.g., guaranteeing one vote per user).

---

## 🔑 Core Features & Security

### 1. Identity & Access Management (IAM)

- **Stateless JWT Authentication:** Issues both short-lived Access Tokens and long-lived Refresh Tokens.
- **Secure Logout (Blacklisting):** Logs out users by placing their active JWT into a Thread-Safe In-Memory Cache (Blacklist) until its natural expiration.
- **Role-Based Access Control (RBAC):** Custom Gin middleware protects administrative routes (e.g., `Promote User`, `Delete Any Question`).
- **Bcrypt Hashing:** Passwords are mathematically salted and hashed before entering the database.

### 2. The Question Engine

- **Advanced Search & Pagination:** Cursor-based pagination (`?cursor=timestamp`) and regex-based keyword search (`?search=linkedlist`) combined with exact tag/company filtering.
- **Atomic Interactions:** Tracks Upvotes, Downvotes, and Views. Uses MongoDB's atomic `$inc` operators to ensure perfect metrics under heavy concurrent load.
- **Anti-Spam Voting System:** Utilizes a dedicated `Vote` Truth Table and a Compound Unique Index to guarantee one vote per user per question.

### 3. Data Integrity

- **Soft Deletes / Anonymization:** Users can be deactivated (`DeletedAt`). Their profile is hidden, but their valuable interview questions remain on the platform as `[Anonymous]`, preventing cascading data loss.

### 4. AI Integration

- **Generative AI:** Integrates with the Google Gemini API to automatically generate study hints for interview questions on demand.

---

## ⚙️ Local Development Setup

### Prerequisites

- [Go 1.21+](https://go.dev/dl/)
- [MongoDB Community Server](https://www.mongodb.com/try/download/community) (Running on port 27017)

### 1. Clone & Install

```bash
git clone https://github.com/PhiliposHailu/prepbase
cd prepbase
go mod tidy
```

### 2. Environment Variables

Create a `.env` file in the root directory:

```env
PORT=8080
JWT_ACCESS_SECRET=your_super_secret_access_key
JWT_REFRESH_SECRET=your_super_secret_refresh_key
DB_URI=mongodb://localhost:27017
DB_NAME=prepbase_db
GEMINI_API_KEY=your_gemini_api_key_here
```

### 3. Run the Server

```bash
go run main.go
```

_The server will start on `http://localhost:8000`._

---

## 🧪 Testing

Unit tests are written using the Table-Driven design pattern and `testify/mock`. Because of the Clean Architecture, tests run in milliseconds without requiring a live MongoDB connection.

```bash
# Generate fresh mocks (if interfaces changed)
go run github.com/vektra/mockery/v2@latest --dir=domain --all --output=mocks --outpkg=mocks

# Run the test suite
go test ./usecase/... -v
```

## 📚 API Documentation

A complete Postman Collection containing all routes, required payloads, and example responses is included in this repository.

1. Open Postman.
2. Click **Import**.
3. Select the `PrepBase_API.postman_collection.json` file located in the `/docs` folder.
