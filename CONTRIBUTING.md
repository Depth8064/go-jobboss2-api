# Contributing

Thanks for contributing.

## Development Setup

1. Install Go 1.26 or newer.
2. Clone the repository.
3. Run tests:

```bash
go test ./...
```

## Quality Gate

Before opening a pull request, run:

```bash
go test ./... -coverprofile=coverage.out -covermode=count
go tool cover -func=coverage.out | tail -n 1
```

## Pull Request Checklist

1. Keep changes focused and scoped.
2. Add or update tests for behavior changes.
3. Ensure tests pass locally.
4. Update README when public usage changes.
5. Use clear commit messages.
