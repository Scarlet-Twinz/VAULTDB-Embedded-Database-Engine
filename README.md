# VAULTDB

VAULTDB is an embedded relational database engine written in Go. The implementation makes core database internals explicit rather than hiding them behind a framework.

## Engine capabilities

- Fixed-size persistent pages and page files
- Length-prefixed record storage
- Buffer pool with pinning, dirty tracking and LRU eviction
- Persistent schema catalog
- SQL parsing for table creation, inserts, selects, updates, deletes and transaction control
- Relational execution operators for scan, filter, projection, joins, grouping, aggregation, sorting and limits
- B+ tree index implementation with ordered lookup, insertion, splitting, persistence support and deletion
- Transaction manager with BEGIN, COMMIT and ROLLBACK records
- Append-only write-ahead logging with LSNs and durable flushes
- WAL transaction analysis for recovery/replay selection
- Lock manager for read/write resource coordination
- Failure-injection primitive for deterministic fault tests
- Thread-safe engine metrics
- CLI / REPL
- Unit tests, integration tests, race detection, vetting, build verification and benchmarks

## Architecture

CLI
 -> SQL
 -> Planner
 -> Execution Operators
 -> Transaction Manager
 -> Buffer Pool
 -> Page Storage
 -> Disk

WAL runs alongside the transaction manager. The catalog persists schema metadata and the index subsystem provides ordered key access.

## Verification

The required CI gate is:

- gofmt
- go vet ./...
- go test ./...
- go test -race ./...
- go build ./...

Additional benchmark and subsystem tests live beside the implementation.

## Scope

VAULTDB is an engineering and learning database implementation. It is not presented as a production replacement for PostgreSQL, MySQL or SQLite.

## Development

go test ./...
go test -race ./...
go vet ./...
go build ./...
go run ./cmd/vaultdb --database ./vaultdb-data

Example SQL:

CREATE TABLE users (id INT, name TEXT);
INSERT INTO users VALUES (1, 'Ada');
SELECT * FROM users;

## License

MIT
