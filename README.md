# Pipedrive Deals API

A simple proxy to the Pipedrive API for managing deals.

## Project Structure

```
pdapi_hometask/
├── go.mod
├── main.go
├── handlers/
│   └── deals.go
├── middleware/
│   └── middleware.go
├── models/
│   └── deal.go
└── utils/
    └── forward.go
```

## Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/yourusername/pdapi_hometask.git
    ```
2. Navigate to the project directory:
    ```bash
    cd pdapi_hometask
    ```
3. Install the dependencies:
    ```bash
    go mod tidy
    ```

## Usage

1. Set the required environment variables:
    ```bash
    export PIPEDRIVE_API_TOKEN=your_api_token
    export PIPEDRIVE_COMPANY_DOMAIN=your_company_domain
    ```
2. Generate Swagger documentation:
    ```bash
    swag init
    ```
3. Run the application:
    ```bash
    go run .
    ```
4. Open Swagger UI:
    Navigate to [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) in your browser to view and test your API endpoints.

## Endpoints

- `GET /deals`: Retrieve all deals.
- `POST /deals`: Create a new deal.
- `PUT /deals/{id}`: Update an existing deal.
- `GET /metrics`: Prometheus metrics endpoint.

## Author
vnahynal


