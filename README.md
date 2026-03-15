# PrepBase API 🚀

PrepBase is a backend REST API built for developers to share, discuss, and study technical interview questions. It is built using **Go (Golang)**, **Gin**, and **MongoDB**, strictly adhering to **Clean Architecture** principles.

## Features

- **Authentication & RBAC:** JWT-based stateless authentication with User and Admin roles.
- **Anti-Spam Voting:** A robust upvote/downvote engine ensuring unique votes per user.
- **AI Integration:** Asynchronous Google Gemini AI integration to generate study hints.
- **Thread-Safe Caching:** In-memory caching using `sync.RWMutex` for high-traffic endpoints.
- **Soft Deletes:** Anonymization of user data preserving database integrity.

## How to Run

1. **Clone the repository.**
2. **Setup MongoDB:** Ensure MongoDB is running locally on port `27017`.
3. **Environment Variables:** Rename `.env.example` to `.env` and insert your JWT secrets and Gemini API Key.
4. **Install Dependencies:**
   ```bash
   go mod tidy
   ```
5. **Run the Server:**
   ```bash
   go run main.go
   ```

## Testing

Unit tests are implemented using `testify/mock` via Table-Driven design.

```bash
go test ./usecase/... -v
```

## API Documentation

A full Postman collection with example requests and responses is located in the `/docs` folder. Import this into Postman to interact with the API.
