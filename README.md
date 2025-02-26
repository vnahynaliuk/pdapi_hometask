# Pipedrive Deals API

A Go application that acts as a proxy to the [Pipedrive API](https://developers.pipedrive.com/docs/api/v1/Deals) for managing deals. The API provides endpoints for retrieving, creating, and updating deals. It also integrates logging, Prometheus metrics, Swagger documentation, and CI/CD (via GitHub Actions).

---

## Table of Contents

1. [Project Structure](#project-structure)  
2. [Features](#features)  
3. [Prerequisites](#prerequisites)  
4. [Installation](#installation)  
5. [Configuration](#configuration)  
6. [Running Locally](#running-locally)  
7. [Docker Setup](#docker-setup)  
8. [Swagger Documentation](#swagger-documentation)  
9. [Prometheus Metrics](#prometheus-metrics)  
10. [Testing](#testing)  
11. [CI/CD](#cicd)  
12. [How It Works](#how-it-works)  
13. [License](#license)

---

## Project Structure

The project is organized as follows:

```
pdapi_hometask/
├── go.mod
├── main.go
├── handlers/
│   ├── deals.go          # Contains endpoint handlers for deals (GET, POST, PUT)
│   └── deals_test.go     # Unit tests for deal handlers
├── middleware/
│   └── middleware.go     # Middleware for logging and Prometheus metrics
├── models/
│   └── deal.go           # Request body models (CreateDeal, UpdateDeal) with example tags
├── utils/
│   └── forward.go        # Utility function to forward HTTP requests
└── docs/                 # Auto-generated Swagger documentation (if available)
```

---

## Features

- **Endpoints:**
  - **GET /deals**: Retrieve all deals.
  - **POST /deals**: Create a new deal.
  - **PUT /deals/{id}**: Update an existing deal.
- **Logging**: Each request is logged with its method, URI, and remote address.
- **Prometheus Metrics**: Metrics (e.g., total requests and request duration) are collected via middleware and available at **GET /metrics**.
- **Swagger Documentation**: Automatically generated API docs available at **GET /swagger/index.html**.
- **Docker Support**: Easily build and run the application inside a container.
- **Testing**: Unit tests are included for endpoint handlers.
- **CI/CD**: GitHub Actions workflows for continuous integration and deployment (setup available in the repository).

---

## Prerequisites

- **Go** 1.24.0 (or higher)
- **Git** (for cloning the repository)
- **Docker** (for containerized deployment)
---

## Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/yourusername/pdapi_hometask.git
   cd pdapi_hometask
   ```

2. **Download Go modules:**

   ```bash
   go mod download
   ```

---

## Configuration

Set the following environment variables (e.g., in your terminal or using a .env file):

- **PIPEDRIVE_API_TOKEN** – Your Pipedrive API token  
- **PIPEDRIVE_COMPANY_DOMAIN** – Your Pipedrive company subdomain (e.g., if your URL is `valeriia-sandbox.pipedrive.com`, then set this value to `valeriia-sandbox`)

Example (macOS/Linux):

```bash
export PIPEDRIVE_API_TOKEN="your_api_key"
export PIPEDRIVE_COMPANY_DOMAIN="your_company_domain"
```

---

## Running Locally

After setting the required environment variables:

```bash
go run .
```

The server listens on port **8080**. You can test endpoints:

- **GET Deals:** [http://localhost:8080/deals](http://localhost:8080/deals)
- **Swagger UI:** [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
- **Prometheus Metrics:** [http://localhost:8080/metrics](http://localhost:8080/metrics)

---

## Docker Setup

### 1. Build the Docker Image

From the project root, run:

```bash
docker build -t pdapi-hometask .
```

### 2. Run the Docker Container

Run the container with the necessary environment variables:

```bash
docker run -p 8080:8080 \
  -e PIPEDRIVE_API_TOKEN="your_api_key" \
  -e PIPEDRIVE_COMPANY_DOMAIN="your_company_domain" \
  pdapi-hometask
```

### 3. Testing in Docker

Open your browser or use curl:

- **GET /deals:** [http://localhost:8080/deals](http://localhost:8080/deals)
- **Swagger UI:** [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
- **Metrics:** [http://localhost:8080/metrics](http://localhost:8080/metrics)

---

## Swagger Documentation

This project uses [swaggo](https://github.com/swaggo/swag) to auto-generate API docs.

1. **Install the swag tool (if not installed):**

   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

2. **Generate the Swagger docs:**

   ```bash
   swag init
   ```

   This creates a **docs/** folder containing `swagger.json` and `swagger.yaml`.

3. **Access Swagger UI:**  
   Run the server and navigate to [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html).

---

## Prometheus Metrics

The application uses Prometheus client libraries to collect metrics via middleware. Metrics available:

- **http_requests_total:** Total number of HTTP requests, labeled by method and endpoint.
- **http_request_duration_seconds:** Histogram for request durations.

Access metrics at [http://localhost:8080/metrics](http://localhost:8080/metrics).

---

## Testing

### Unit Tests

Unit tests for endpoint handlers are located in `handlers/deals_test.go`. To run the tests:

```bash
go test ./...
```

### Manual Testing

You can use curl, Postman, or Swagger UI to test endpoints. For example:

- **GET deals:**

  ```bash
  curl http://localhost:8080/deals
  ```

- **POST deal:**  
  Use the following JSON body as an example:
  
  ```json
  {
    "title": "My New Deal",
    "value": "5000",
    "label": [1, 2],
    "currency": "USD",
    "visible_to": "3"
  }
  ```

- **PUT deal:**  
  Use a similar JSON body to update an existing deal (replace `{id}` with the actual deal ID).

---

## CI/CD

### GitHub Actions

The repository includes GitHub Actions workflows:

- **CI Workflow (`.github/workflows/ci.yml`):**  
  - Runs on every push and pull request.
  - Installs dependencies, runs tests, and checks code quality with `go vet`.

- **CD Workflow (`.github/workflows/cd.yml`):**  
  - Triggered when a pull request is merged into the `main` branch.
  - Performs deployment steps (or logs a deployment message).

Check the **Actions** tab in your GitHub repository for details.

---

## How It Works

- **Main Application:**  
  The `main.go` file sets up a router using Gorilla Mux, applies logging and metrics middleware, and registers endpoints for deals, metrics, and Swagger UI.

- **Handlers:**  
  The `handlers/deals.go` file defines endpoint functions for GET, POST, and PUT requests. It reads environment variables, constructs URLs to the Pipedrive API, forwards requests via a utility function, and returns responses.

- **Models:**  
  The `models/deal.go` file defines the request body structures (`CreateDeal` and `UpdateDeal`) with example tags so that Swagger UI displays an example JSON. For instance, the example for creating a deal is:

  ```json
  {
    "title": "My New Deal",
    "value": "5000",
    "label": [1, 2],
    "currency": "USD",
    "visible_to": "3"
  }
  ```

- **Middleware:**  
  The `middleware/middleware.go` file contains functions to log each request and collect Prometheus metrics such as request counts and durations.

- **Utils:**  
  The `utils/forward.go` file provides a helper function (`ForwardRequest`) to create and send HTTP requests.

- **Swagger Docs:**  
  Documentation is auto-generated by swaggo from comments in the code. Run `swag init` to generate/update the docs, which are then served via the `/swagger/index.html` endpoint.

- **Docker:**  
  The Dockerfile builds the project into an image and runs it. Environment variables are passed to the container, ensuring that the application connects to the correct Pipedrive account. See the [Docker Setup](#docker-setup) section for details.

---

## Author
[vnahynaliuk](https://github.com/vnahynaliuk)

