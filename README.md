# VAULTDB

**An embedded relational database engine written in Go, built from the storage layer upward.**

VAULTDB is a from-scratch database-engineering project that makes the internals of a relational database explicit: pages, records, buffering, catalogs, indexes, SQL parsing, execution operators, transactions, WAL, recovery analysis, concurrency, fault injection, and observability.

The project is intentionally educational and engineering-focused. It is not intended to replace PostgreSQL, MySQL, or SQLite.

---

## Why VAULTDB exists

Most application projects use a database without exposing what happens underneath it.

VAULTDB explores the other side of that boundary:

> **What does it take to build a relational database engine from the ground up?**

The project is organized around the path from a SQL statement to persistent storage:

**SQL → planning/execution → transactions → buffering → pages → disk**

Along the way, the repository contains independently testable database subsystems rather than hiding the important mechanics behind an external database framework.

---

## What is implemented

### Storage

- Fixed-size 4096-byte pages
- Persistent page files
- Page encoding and decoding
- Length-prefixed records
- Buffer pool
- Page pinning
- Dirty-page tracking
- LRU eviction
- Explicit flushing and synchronization

### Catalog

- Persistent schema catalog
- Table definitions
- Column metadata
- Index metadata
- Table and index lookup

### SQL and execution

The SQL layer currently supports:

- `CREATE TABLE`
- `INSERT`
- `SELECT`
- `UPDATE`
- `DELETE`
- `BEGIN`
- `COMMIT`
- `ROLLBACK`
- `SHOW TABLES`

The repository also contains independently tested relational execution operators for:

- scans
- filtering
- projection
- joins
- grouping
- aggregation
- sorting
- limits

The planner package models corresponding relational plan types. These subsystems are deliberately kept explicit and testable rather than being presented as a fully cost-based SQL optimizer.

### Indexing

VAULTDB includes a B+ tree implementation with:

- ordered string-key lookup
- insertion
- node splitting
- leaf links
- persistence support
- deletion/rebuild handling
- dedicated tests
- benchmarks

### Transactions and WAL

The transaction subsystem provides:

- transaction IDs
- transaction state tracking
- BEGIN / COMMIT / ROLLBACK records
- append-only write-ahead logging
- log sequence numbers
- durable WAL synchronization

The recovery package analyzes WAL transaction state and identifies committed records for replay selection.

This is an explicit recovery foundation, not a claim of production-grade crash recovery or full MVCC.

### Concurrency

The repository includes a read/write lock manager with conflict detection and dedicated concurrency tests.

### Fault injection

A deterministic failure-injection primitive is included so failure scenarios can be tested without relying on random failures.

### Metrics

Thread-safe metrics track database activity such as:

- pages read
- pages written
- buffer hits and misses
- index lookups
- table scans
- transactions
- queries
- query execution time

### CLI

VAULTDB includes a command-line interface and REPL.

You can provide a database directory with:

```bash
go run ./cmd/vaultdb --database ./vaultdb-data
```

or execute a statement directly:

```bash
go run ./cmd/vaultdb --database ./vaultdb-data --execute "SHOW TABLES"
```

---

## Architecture

The intended engine flow is:

```text
                         Client / CLI
                              │
                              ▼
                         SQL Parser
                              │
                              ▼
                       Query Planning
                              │
                              ▼
                    Execution Operators
                              │
                              ▼
                     Transaction Manager
                       │      │       │
                       │      │       └── WAL
                       │      │
                       │      └────────── Index Manager
                       │
                       ▼
                      Buffer Pool
                              │
                              ▼
                         Page Storage
                              │
                              ▼
                             Disk
```

Supporting subsystems include:

- persistent catalog metadata
- B+ tree indexing
- read/write locking
- recovery analysis
- deterministic fault injection
- engine metrics

The repository keeps these layers separated so individual database mechanisms can be tested without requiring the complete engine path for every test.

---

## Example session

Create a table:

```sql
CREATE TABLE users (id INT, name TEXT);
```

Insert records:

```sql
INSERT INTO users VALUES (1, 'Ada');
INSERT INTO users VALUES (2, 'Grace');
```

Read them:

```sql
SELECT * FROM users;
```

Update a record:

```sql
UPDATE users SET name = 'Ada Lovelace' WHERE id = 1;
```

Delete a record:

```sql
DELETE FROM users WHERE id = 2;
```

Inspect the catalog:

```sql
SHOW TABLES;
```

Transaction control is also available:

```sql
BEGIN;
COMMIT;
```

and:

```sql
BEGIN;
ROLLBACK;
```

---

## How to shake it out

If you want to review VAULTDB like an engineer rather than simply reading the source, start with the repository's automated gate:

```bash
go test ./...
go test -race ./...
go vet ./...
go build ./...
```

Then run the database against a persistent directory:

```bash
go run ./cmd/vaultdb --database ./vaultdb-data
```

Exercise table creation, inserts, reads, updates, deletes, and transaction commands.

Then exit and reopen VAULTDB using the **same** `./vaultdb-data` directory. Verify that persisted state is still available.

For deeper subsystem testing, inspect the package tests under:

```text
internal/storage/
internal/catalog/
internal/sql/
internal/engine/
internal/index/
internal/transaction/
internal/wal/
internal/recovery/
internal/concurrency/
internal/execution/
internal/faults/
internal/metrics/
```

The race-enabled test run is particularly important for checking the concurrency-sensitive code paths.

---

## Verification

The repository's CI gate runs:

1. **gofmt**
2. **go vet ./...**
3. **go test ./...**
4. **go test -race ./...**
5. **go build ./...**

The repository also contains subsystem tests and a B+ tree benchmark.

The latest completed CI verification for the completed implementation passed all required gate stages.

---

## Testing philosophy

VAULTDB treats database internals as independently testable engineering components.

Tests cover areas including:

- page encoding and persistence
- records
- buffer-pool behavior
- catalog persistence
- SQL parsing
- engine persistence
- B+ tree operations
- transactions
- WAL behavior
- recovery analysis
- relational operators
- locking
- fault injection
- metrics

Race detection is included in the required CI gate so concurrency issues are not checked only by ordinary unit tests.

---

## Repository structure

```text
VAULTDB-Embedded-Database-Engine/
├── cmd/
│   └── vaultdb/              # CLI / REPL entry point
├── internal/
│   ├── catalog/              # persistent schema catalog
│   ├── concurrency/          # read/write lock manager
│   ├── engine/               # database engine orchestration
│   ├── execution/            # relational execution operators
│   ├── faults/               # deterministic failure injection
│   ├── index/                 # B+ tree indexes
│   ├── metrics/               # thread-safe metrics
│   ├── planner/               # relational plan structures
│   ├── recovery/              # WAL transaction analysis
│   ├── sql/                   # SQL parser and statements
│   ├── storage/               # pages, records and buffer pool
│   ├── transaction/           # transaction lifecycle
│   └── wal/                   # write-ahead log
├── .github/
│   └── workflows/
│       └── ci.yml             # formatting, vet, tests, race, build
├── go.mod
├── LICENSE
└── README.md
```

---

## Technology

| Area | Technology |
|---|---|
| Language | Go |
| Database model | Embedded relational |
| Storage | Fixed-size pages and page files |
| Index | B+ tree |
| Logging | Write-ahead log |
| Concurrency | Read/write locking |
| Interface | CLI / REPL |
| Testing | Go tests + race detector |
| CI | GitHub Actions |
| Dependencies | Intentionally minimal |

---

## Run locally

Requirements:

- Go 1.23+
- Git
- A terminal

Clone the repository:

```bash
git clone https://github.com/Scarlet-Twinz/VAULTDB-Embedded-Database-Engine.git
cd VAULTDB-Embedded-Database-Engine
```

Run the test suite:

```bash
go test ./...
```

Run the race detector:

```bash
go test -race ./...
```

Run static analysis:

```bash
go vet ./...
```

Build everything:

```bash
go build ./...
```

Start the REPL:

```bash
go run ./cmd/vaultdb --database ./vaultdb-data
```

---

## Scope and engineering boundaries

VAULTDB is deliberately a database-engineering project, not a claim that a small from-scratch implementation should replace mature production databases.

The repository demonstrates core mechanisms and their tests while keeping the scope understandable.

Some subsystems currently exist as explicit, independently tested building blocks rather than being presented as fully integrated production features. In particular:

- the planner is not a cost-based optimizer;
- the relational operators are not all exposed through the current SQL surface;
- the lock manager is an explicit concurrency subsystem rather than a claim of complete database isolation semantics;
- WAL recovery currently performs transaction-state analysis rather than full production REDO/UNDO crash recovery;
- fault injection and metrics provide reusable infrastructure but are not presented as comprehensive instrumentation of every engine path.

Those boundaries are intentional. They make the README match the implementation instead of overstating what the engine does.

---

## Design principles

### Build from the storage layer upward

The project starts with pages and records rather than beginning with a UI and hiding persistence behind an external database.

### Make database internals visible

Important mechanisms live in their own packages so they can be inspected, reasoned about, and tested independently.

### Persistence matters

A database is more than an in-memory collection. VAULTDB exercises page files, catalogs, WAL, indexes, and reopen/persistence behavior.

### Test the mechanisms

The project uses unit tests, integration tests, race detection, static analysis, builds, and benchmarks rather than relying only on manual demonstrations.

### Be explicit about boundaries

VAULTDB distinguishes implemented mechanisms from future production-grade concerns instead of treating a prototype as a complete PostgreSQL replacement.

---

## Current status

**Complete engineering milestone — implemented, tested, race-checked, vetted, built, and verified through GitHub Actions.**

The current repository contains the storage, catalog, SQL, execution, indexing, transaction, WAL, recovery-analysis, concurrency, fault-injection, metrics, CLI, testing, and CI foundations described above.

The required CI gate is green for the completed implementation.

---

## Project links

- **Repository:** https://github.com/Scarlet-Twinz/VAULTDB-Embedded-Database-Engine
- **GitHub:** https://github.com/Scarlet-Twinz

---

## Author

**Anthony Emmanuella Mmasinachi**

Full-stack and systems-focused developer building projects across web applications, backend systems, SaaS architecture, automation, AI integration, and practical software engineering.

---

## License

MIT License.

See [LICENSE](LICENSE) for the full license text.
