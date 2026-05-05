# Invoice View

A Go-based invoice management service built with Fiber web framework. This application provides RESTful API endpoints for creating invoices and viewing them as rendered HTML pages.

## Features

- **RESTful API**: Create invoices via POST endpoint with API key authentication
- **HTML Invoice View**: Render invoices as beautifully formatted HTML pages
- **Template System**: Support for multiple invoice templates
- **Rate Limiting**: Built-in rate limiting middleware to protect endpoints
- **SQLite Database**: Lightweight storage with WAL mode for better concurrency
- **Currency Formatting**: Support for various currency formats
- **Status Management**: Track invoice status (PAID, PENDING, EXPIRED, CANCELLED)

## Project Structure

```
invoice-view/
├── database/          # Database connection and migrations
├── model/             # Data models and structs
├── repository/        # Data access layer
├── router/            # HTTP routes and handlers
│   └── middleware/    # Custom middleware (rate limiter, API key auth)
├── service/           # Business logic layer
├── template/          # HTML invoice templates
├── util/              # Utility functions
├── main.go            # Application entry point
├── go.mod             # Go module definition
└── .env.sample        # Environment variables template
```

## Requirements

- Go 1.23.12 or higher
- SQLite3

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd invoice-view
```

2. Install dependencies:
```bash
go mod download
```

3. Configure environment variables:
```bash
cp .env.sample .env
```

Edit `.env` with your preferred settings:
```env
APP_PORT=2121
LIMITER_INVOICE_MAX=60
LIMITER_INVOICE_DURATION=1
API_SECRET=rahasiabanget
```

## Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_PORT` | Server port | `3000` |
| `LIMITER_INVOICE_MAX` | Maximum requests per duration | `60` |
| `LIMITER_INVOICE_DURATION` | Rate limit duration (minutes) | `1` |
| `API_SECRET` | API key for authentication | - |

## Running the Application

```bash
go run main.go
```

The server will start on the configured port (default: `:3000`).

## API Endpoints

### Create Invoice

**POST** `/v1/invoice`

Creates a new invoice in the system.

**Headers:**
- `X-API-Key`: Your API secret key
- `Content-Type`: `application/json`

**Request Body:**
```json
{
  "id": "inv_001",
  "invoice_number": "INV-2024-001",
  "currency": "IDR",
  "amount": 1500000,
  "status": "PENDING",
  "issued_at": "2024-01-15T10:00:00Z",
  "expired_at": "2024-02-15T10:00:00Z",
  "customer": {
    "id": "cust_001",
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "+628123456789",
    "address": "Jakarta, Indonesia"
  },
  "items": [
    {
      "sku": "PROD-001",
      "name": "Product Name",
      "qty": 2,
      "price": 500000,
      "total": 1000000
    }
  ],
  "summary": {
    "subtotal": 1000000,
    "tax": 100000,
    "discount": 0,
    "grand_total": 1500000
  },
  "metadata": {
    "notes": "Thank you for your business"
  },
  "config": {
    "template": "invoice",
    "orientation": "portrait"
  }
}
```

**Response:** `201 Created`
```json
{
  "id": "inv_001",
  "invoice_number": "INV-2024-001",
  ...
}
```

### View Invoice HTML

**GET** `/invoice/:id`

Renders an invoice as an HTML page.

**Response:** `text/html`

Example: `GET /invoice/inv_001`

## Templates

The application supports multiple invoice templates located in the `template/` directory:

- `invoice.html` - Default template
- `invoice_2.html` - Alternative template

Templates can be selected via the `config.template` field when creating an invoice.

## Middleware

### Rate Limiter
Applied to `/invoice/:id` endpoint to prevent abuse. Configurable via environment variables.

### API Key Authentication
Required for `/v1/invoice` POST endpoint. Validate using `X-API-Key` header.

## Development

### Building

```bash
go build -o invoice-view main.go
```

### Testing

Run tests with:
```bash
go test ./...
```

## Technologies Used

- **[Fiber](https://gofiber.io/)** - Fast, Express-like web framework for Go
- **[SQLite3](https://github.com/mattn/go-sqlite3)** - Embedded database driver
- **[godotenv](https://github.com/joho/godotenv)** - Environment variable management
- **[msgp](https://github.com/tinylib/msgp)** - MessagePack serialization
- **[brotli](https://github.com/andybalholm/brotli)** - Compression support

## License

This project is open source and available under the MIT License.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
