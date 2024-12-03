# Music Library API

Stack: Golang + Postgres + Docker + Swagger

## How to run
Configuration:
- [dev config](./config/dev.yml) - for local running and testing
- [deploy config](./config/deploy.yml) - for deploying in a Docker container

Run locally:
```bash
go run cmd/main.go --config=config/dev.yml
```

Run in Docker:
```bash
docker compose up --build
```

## How Pagination Works

This project uses cursor pagination. Compared to indexed pagination, this method is faster, reduces database load, and decreases response time.

The `OFFSET` instruction is costly because, in a large database, it requires iterating through many records, increasing response time. This is mitigated by using `ORDER BY id > cursor_id`.