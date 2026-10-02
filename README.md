# TX Technical Challenge

Small Go service for finding profitable currency arbitrage routes from a set of market exchange rates.

The project exposes a gRPC API, converts incoming proto markets into domain models, builds an exchange graph, enumerates closed routes, and returns routes whose combined exchange rate is profitable.

## How To Read The Codebase

Start from the outside and move inward:

1. `main.go`
   - Application entrypoint.
   - Delegates startup to the `cmd` package.

2. `cmd/start.go`
   - Selects the command to run.
   - With no arguments, starts the gRPC server.
   - Supports:
     - `serve`
     - `request`

3. `cmd/serve.go`
   - Opens the TCP listener.
   - Creates the gRPC server.
   - Registers the API implementation from `internal/server`.

4. `internal/server/grpc.go`
   - gRPC handler layer.
   - Validates and converts protobuf input into arbitrage domain objects.
   - Calls the arbitrage engine.
   - Converts routes back into protobuf response objects.

5. `internal/arbitrage`
   - Core business logic.
   - This package does not depend on gRPC or protobuf.

## Running Locally

Run the gRPC server

### Locally

```sh
go run . serve
```

### Docker

Build the image:

```sh
make docker-build
```

Run it:

```sh
make docker-run
```

### Send a Request

You could use external GRPC tools (e.g. Postman) or send the built-in sample request:

```sh
go run . request
```

Message example:

```json
{
    "markets": [
        {
            "base_currency": "EUR",
            "quote_currency": "USD",
            "exchange_rate": "1.1586"
        },
        {
            "base_currency": "EUR",
            "quote_currency": "GBP",
            "exchange_rate": "0.6849"
        },
        {
            "base_currency": "USD",
            "quote_currency": "GBP",
            "exchange_rate": "0.5904"
        }
    ]
}
```
