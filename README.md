# Sensor Telemetry Ingestion Service

This Go service models the ingestion path for industrial sensor readings. It includes snapshot aggregation, streamed batch intake, downstream retry classification, profile enrichment, cancellation-aware polling, binary frame ownership, capture processing, delivery state tracking, background dispatch, and durable event publication.

## Layout

- `cmd/server`: runnable service entrypoint
- `internal/model`: shared telemetry types
- `internal/store` and `internal/aggregate`: concurrent reading snapshots
- `internal/source` and `internal/batch`: streamed batch collection
- `internal/adapter`, `internal/journal`, and `internal/transmit`: downstream error handling
- `internal/config` and `internal/enrich`: sensor profile loading
- `internal/poller` and `internal/backoff`: cancellation-aware polling
- `internal/decoder`, `internal/cache`, and `internal/exporter`: binary frame ownership
- `internal/capture` and `internal/transaction`: capture resource processing
- `internal/state` and `internal/delivery`: retry state and idempotency
- `internal/dispatcher` and `internal/scheduler`: background work lifecycle
- `internal/outbox`, `internal/audit`, and `internal/publish`: durable event delivery

## Run

```bash
go run ./cmd/server
```

## Verify

```bash
go build ./...
go test ./...
```
