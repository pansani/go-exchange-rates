# Go Exchange Rate 

This project is a simple Go application for fetching USD to BRL exchange rates from an external API, storing the rates in a SQLite database, and exposing an endpoint to retrieve the latest rate.

---

## Features

- **API Integration**: Fetches USD to BRL exchange rates from `https://economia.awesomeapi.com.br/json/last/USD-BRL`.
- **SQLite Database**: Stores the exchange rates locally.
- **REST Endpoint**: Provides an endpoint (`/cotacao`) to retrieve the latest rate.
- **Unit Tests**: Includes comprehensive tests for database operations and API handlers.

---

## Project Structure

```
client-server-full-cycle/
├── client/
│   └── client.go          # Placeholder for potential client-side logic
├── server/
│   ├── server.go          # Main server logic and `/cotacao` handler
│   ├── server_test.go     # Unit tests for server.go
├── app.db                 # SQLite database (created during runtime)
├── cotacao.txt            # Sample output file
├── go.mod                 # Module dependencies
├── go.sum                 # Dependency checksums
└── README.md              # Project documentation
```

---

## Prerequisites

- **Go**: Install the latest version of Go from [golang.org](https://golang.org/dl/).
- **SQLite3**: Ensure SQLite3 is installed on your system.

---

## Installation

1. Clone this repository:
   ```bash
   git clone https://github.com/your-username/client-server-full-cycle.git
   cd client-server-full-cycle
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

---

## Usage

1. Start the server:
   ```bash
   go run server/server.go
   ```
   The server will start on `http://localhost:8080`.

2. Fetch the latest USD to BRL exchange rate:
   ```bash
   curl http://localhost:8080/cotacao
   ```
   The response will include the latest rate, e.g.:
   ```
   Dolar: 5.42
   ```

---

## Testing

Run unit tests with:
```bash
go test ./...
```

### Included Tests:
- **Database Operations**: Validates table creation, insertion, and data retrieval.
- **API Handler**: Verifies `/cotacao` endpoint behavior with a mocked external API.




