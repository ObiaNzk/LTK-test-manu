# LTK-test-manu

Events API service built with Go and PostgreSQL.

## Prerequisites

- [Docker Desktop](https://www.docker.com/products/docker-desktop/) - Required for running PostgreSQL
- Go - For running the application (was made with 1.25 but should work with 1.21+)

---

## Quick Start

### 1. Setup Database

```bash
# Pull Postgres container
make postgres-pull
# Start Postgres container
make postgres-start
# Create Events database and run migrations
make db-setup
-- 
```

### 2. Run Application

```bash
make run
```

**Server runs on `http://localhost:8080`**

## API Endpoints

### POST /events

Creates a new event in the database.

**Request Body:**

```json
{
  "title": "pepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepito",
  "description": "hire me, maybe",
  "start_time": "2025-12-01T09:00:00Z",
  "end_time": "2025-12-01T10:00:00Z"
}
```

**Notes:**

- Title must be more than 100 characters.
- start and endtime should have the time.time go format

**Success Response (201 Created):**

```json
{
  "id": "e4f5c6d7-8e9f-4a5b-9c8d-7e6f5a4b3c2d",
  "title": "pepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepito",
  "description": "hire me, maybe",
  "start_time": "2025-12-01T09:00:00Z",
  "end_time": "2025-12-01T10:00:00Z",
  "created_at": "2025-11-27T10:30:00Z"
}
```

**Error Responses:**

- `400 Bad Request` - Invalid input
- `500 Internal Server Error` - Database or server error

---

### GET /events

Returns all events ordered by start time (ascending).

**Success Response (200 OK):**

```json
[
  {
    "id": "e4f5c6d7-8e9f-4a5b-9c8d-7e6f5a4b3c2c",
    "title": "pepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepito",
    "description": "hire me, maybe",
    "start_time": "2025-12-01T09:00:00Z",
    "end_time": "2025-12-01T10:00:00Z",
    "created_at": "2025-11-27T10:30:00Z"
  },
  {
    "id": "e4f5c6d7-8e9f-4a5b-9c8d-7e6f5a4b3c2d",
    "title": "asdpepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepito",
    "description": "hire me 2, maybe",
    "start_time": "2025-12-02T09:00:00Z",
    "end_time": "2025-12-02T10:00:00Z",
    "created_at": "2025-11-28T10:30:00Z"
  }
]
```

### GET /events/{id}

Returns a specific event by ID.

**Success Response (200 OK):**

```json
  {
  "id": "e4f5c6d7-8e9f-4a5b-9c8d-7e6f5a4b3c2d",
  "title": "asdpepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepito",
  "description": "hire me 2, maybe",
  "start_time": "2025-12-02T09:00:00Z",
  "end_time": "2025-12-02T10:00:00Z",
  "created_at": "2025-11-28T10:30:00Z"
}
```

**Error Responses:**

- `404 Not Found` - Event not found
- `500 Internal Server Error` - Database or server error

---

### Test with Postman / curl

```sql
curl --location 'localhost:8080/events' \
--header 'Content-Type: application/json' \
--data '{
"title": "asdpepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepitopepito",
"description": "pepito",
"start_time": "2025-12-15T14:00:00Z",
"end_time": "2025-12-15T15:30:00Z"
}
```
Or download https://www.postman.com/downloads/ and use the example request to test the endpoints

## Database Schema


```sql
CREATE TABLE events
(
    id          VARCHAR(36) PRIMARY KEY,
    title       TEXT      NOT NULL,
    description TEXT,
    start_time  TIMESTAMP NOT NULL,
    end_time    TIMESTAMP NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);
```


### Project Structure
This follows the https://www.ardanlabs.com/blog/2017/02/package-oriented-design.html oriented package design with a few
personal modifications that i have learnt working on go

```
├── cmd/api/              
├── internal/             
│   ├── migrations/       
│   ├── platform/         
│   └── service.go
│   └── storage.go  
│   └── errors.go
│   └── models.go              
└── Makefile             
```

