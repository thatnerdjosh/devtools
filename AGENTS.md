# Repository Guidelines

## Project Structure & Module Organization
- `cmdline/go/`: Go module and Makefile for CLI tools. Build output is written as `cmdline/go/build` (directory) and/or `cmdline/go/build` (binary name via `-o`).
- `docker/doctl/` and `docker/kubectl/`: Docker Compose files for `doctl` and `kubectl` helper containers.
- `.env` is ignored by Git; use it for local-only environment variables.

## Build, Test, and Development Commands
- Build Go tools: `make -C cmdline/go build` — creates the `build` directory and compiles packages.
- Clean artifacts: `make -C cmdline/go clean` — removes the `build` directory.
- Direct Go build: `cd cmdline/go && go build ./...` — compile packages without Make.
- Run Kubernetes helper: `docker compose -f docker/kubectl/docker-compose.yml run --rm kubectl`.
- Run DigitalOcean helper: `docker compose -f docker/doctl/docker-compose.yml run --rm doctl` (requires `DIGITALOCEAN_ACCESS_TOKEN`).

## Coding Style & Naming Conventions
- Go code must be formatted with `gofmt`/`goimports`; CI and reviews assume canonical Go style.
- Naming: exported identifiers `PascalCase`, unexported `camelCase`. Package names are short, lowercase, no underscores.
- YAML (Compose): 2-space indents; keep services minimal and explicit.

## Testing Guidelines
- Use Go’s standard testing: files end with `_test.go`, tests named `TestXxx`.
- Run all tests: `cd cmdline/go && go test ./...`.
- Prefer table-driven tests; keep unit tests fast and hermetic.

## Commit & Pull Request Guidelines
- Commit messages: `type(scope): summary` (e.g., `build(go): add Makefile target`).
- PRs include: concise description, rationale, before/after notes, and any required env/setup steps. Link related issues.
- Keep changes small and focused; update docs when behavior or commands change.

## Security & Configuration Tips
- Do not commit secrets. Provide tokens via environment (e.g., `DIGITALOCEAN_ACCESS_TOKEN`) or local `.env` only.
- Validate cloud and cluster contexts before running helper commands.
