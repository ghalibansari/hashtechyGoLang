# Getting Started with Hashtechy API

## Table of Contents

1. [Project Overview](#project-overview)
2. [Prerequisites](#prerequisites)
3. [Installation](#installation)
4. [Project Structure](#project-structure)
5. [Running the Application](#running-the-application)
6. [Development Tools](#development-tools)
7. [Real-time Features](#real-time-features)

## Project Overview

Hashtechy API is a secure user management system that includes:

- 🔒 User data storage with encryption
- ⚡ High-performance Redis caching
- 📨 Asynchronous CSV processing via RabbitMQ
- 🛡️ SSL/TLS security for all services
- 📚 Interactive API documentation (Swagger)
- 🔄 Real-time user updates via SSE

## Prerequisites

Before starting, ensure you have:

### Required Software

- Docker (latest version)
- Docker Compose (latest version)
- Go 1.23.5 or higher
- Visual Studio Code (recommended)
- Web browser with SSE support

### Required Knowledge

- Basic Go programming
- Basic Docker concepts
- REST API fundamentals
- Understanding of Server-Sent Events (SSE)

## Installation

1. **Clone the Repository**

```bash
git clone <repository-url>
cd hashtechy
```

1. **Build and Start Services**

```bash
# Start all services
docker compose up -d

# Verify services are running
docker compose ps
```

## Project Structure

```plaintext
hashtechy/
├── src/                 # Source code
│   ├── encryption/      # Data encryption logic
│   ├── errors/          # Error handling
│   ├── logger/          # Logging system
│   ├── postgres/        # Database operations
│   ├── rabbitmq/        # Message queue handlers
│   ├── redis/           # Cache management
│   ├── server/          # HTTP server & routing
│   └── user/            # User business logic
├── certs/               # SSL certificates
│   ├── ca.crt           # Root CA certificate
│   ├── postgres.crt     # PostgreSQL certificate
│   ├── redis.crt        # Redis certificate
│   └── rabbitmq.crt     # RabbitMQ certificate
├── docs/                # Documentation
├── static/              # Static files
├── docker compose.yml   # Service configuration
└── main.go              # Application entry point
```

## Running the Application

### 1. Start the Services

```bash
docker compose up -d
```

### 2. Verify Services

```bash
# Check service status
docker compose ps

# Check logs
docker compose logs -f app
```

### 3. Access Tools

| Tool                | URL                           | Purpose                  |
| ------------------- | ----------------------------- | ------------------------ |
| Real-time Dashboard | http://localhost:3000/        | SSE User Monitoring      |
| Swagger UI          | http://localhost:3000/swagger | API Documentation        |
| Redis Insight       | http://localhost:5540         | Cache Monitoring         |
| RabbitMQ Console    | http://localhost:15672        | Message Queue Management |

## Development Tools

### 1. API Testing

- **Swagger UI:** Interactive API documentation and testing
- **Postman:** API testing and collection management

### 2. Monitoring

- **Redis Insight:** Cache performance monitoring
- **RabbitMQ Management Console:** Queue monitoring

### 3. Debugging

```bash
# View real-time logs
docker compose logs -f app

# Access container shell
docker compose exec app sh
```

### 4. Common Commands

#### Service Management

```bash
# Restart a service
docker compose restart [service_name]

# View service logs
docker compose logs -f [service_name]

# Stop all services
docker compose down
```

### 5. Swagger Documentation Generation

To generate and update API documentation:

```bash
swag init -g src/server/server.go
```

### 6. Real-time Features

#### SSE Dashboard

The real-time dashboard at http://localhost:3000/ provides:

- Live user count
- Average age statistics
- Real-time user table updates
- Automatic reconnection on connection loss
- Updates every 3 seconds

#### Using SSE Endpoint

The SSE endpoint is available at /events and provides real-time user updates:

```ts
const evtSource = new EventSource("http://localhost:3000/events");
evtSource.addEventListener("userUpdate", (event) => {
  const data = JSON.parse(event.data);
  console.log("Received update:", data);
};
```
