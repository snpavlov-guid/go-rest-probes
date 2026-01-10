# Console Application

A simple Go console application demonstrating basic command-line functionality.

## Building

```bash
go build -o console.exe
```

## Running

```bash
go run main.go
```

Or after building:

```bash
./console.exe
```

## Usage

```bash
console [options] [arguments]

Options:
  -v    Enable verbose output
  -h    Show help message
```

## Examples

```bash
# Basic usage
console

# With arguments
console arg1 arg2 arg3

# Verbose mode
console -v arg1 arg2
```

## API Tests

This project includes comprehensive tests for the aircraft API endpoints.

### Running API Tests

```bash
# Run all tests
go test -v

# Run only API tests
go test -v -run TestGetAircrafts

# Run with custom base URL
go test -v -base-url=http://localhost:8080

# Run with base URL from environment variable
AIRCRAFT_API_BASE_URL=http://localhost:8080 go test -v

# Run authenticated tests (requires JWT token)
AIRCRAFT_API_TOKEN=your-jwt-token go test -v -run TestAuthenticatedEndpoint
```

### Configurable Options

- **Base URL**: Set via `-base-url` flag or `AIRCRAFT_API_BASE_URL` environment variable (default: `http://localhost:8080`)
- **JWT Token**: Set via `AIRCRAFT_API_TOKEN` environment variable for authenticated tests

### Test Coverage

The test suite covers all endpoints from `app_aircraft`:
- Basic endpoints: `/`, `/:text`, `/version`, `/profile`
- Aircraft endpoints: `/api/v1/aircrafts` (GET, POST create/update/delete, DELETE)
- Airport endpoints: `/api/v1/airports` (GET, POST create/update/delete, DELETE)