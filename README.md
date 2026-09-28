# VAULTDB

**VAULTDB is a serious embedded relational database engine written in Go.**

It is built to make database internals explicit: fixed-size pages, durable storage, a buffer pool, a persistent catalog, SQL parsing, execution, transactions, and write-ahead logging. The project is intentionally an engine rather than a web application.

## Current engine slice

- Fixed 4 KiB pages and page headers
- Persistent page-file storage
- Record packing inside pages
- Buffer pool with pinning, dirty tracking and LRU eviction
- Persistent catalog for tables and columns
- SQL parser for `CREATE TABLE`, `INSERT`, `SELECT`, `UPDATE`, `DELETE`, transaction control and `SHOW TABLES`
- Table scans, filtering and projection
- Atomic catalog replacement on save
- Append-only WAL with LSNs and transaction records
- Explicit transaction manager with `BEGIN`, `COMMIT`, `ROLLBACK`
- CLI / REPL
- Unit, integration and race-tested foundation

## Architecture

```text
CLI
 ↓
SQL Parser → AST
 ↓
Execution Layer
 ↓
Transaction Manager ─── WAL
 ↓
Buffer Pool
 ↓
Page Storage
 ↓
Disk
 ↑
Persistent Catalog
```

## Roadmap

The engine is being built vertically. Planned subsystems include B+ tree indexes, stronger record/page layouts, a cost-aware planner, joins and aggregation, crash recovery/REDO, locking and concurrency control, failure injection, fuzzing, and benchmarks. Features are added only with executable tests and explicit persistence/recovery semantics.

## Scope

VAULTDB is an engineering and learning database implementation. It is **not** presented as a production replacement for PostgreSQL, MySQL or SQLite.

## Development

```bash
go test ./...
go test -race ./...
go vet ./...
go build ./...
go run ./cmd/vaultdb --database ./vaultdb-data
```

Example:

```sql
CREATE TABLE users (id INT, name TEXT);
INSERT INTO users VALUES (1, 'Ada');
SELECT * FROM users;
```

## License

MIT
