# Next.js + Golang (Fiber v3) Project Template

Reusable project template combining a **Next.js (`webui/`)** frontend with a **Golang Fiber v3 (`src/`)** backend in a single deployable binary.

## Requirements
- **Go 1.27.4+** with Fiber **v3.7.0**
- **Bun 1.4.5+**
- **Make**
- Frontend baseline: Next.js **16.4.6**, React **19.3.1**, TypeScript **6.0.3**, Tailwind CSS **4.3.5**, and ESLint **10.15.0**

## Architecture
- **Development (`make run`)**: Runs `bun dev` (`:3000`) and Go Fiber (`:8082`) behind a unified Go reverse proxy on `:8080` (`PORT`). `/api/v1/*` requests route to Fiber; all other requests (including Next.js HMR WebSockets) route to Next.js.
- **Production (`make build`)**: Builds the Next.js static export (`webui/out`), copies it to `src/out`, and embeds it via `//go:embed all:out` into a single standalone binary at `build/app`.

## Commands
```bash
make run     # Run Next.js + Go dev server with reverse proxy on :8080
make test    # Run Go unit tests
make build   # Build Next.js static export + single Go binary (build/app)
make clean   # Clean build artifacts
make rename NAME=github.com/you/new_project  # Rename the Go module + all imports
```

## Starting a New Project From This Template
The Go module path is `github.com/hylmithecoder/templatenextjs_golangserver`. Rename it once after creating your project.

**Option A: `git clone` + `make rename`**
```bash
git clone https://github.com/hylmithecoder/templatenextjs_golangserver new_project
cd new_project
make rename NAME=github.com/you/new_project   # or simply NAME=new_project
rm -rf .git && git init                       # optional: start a fresh history
```

**Option B: [`gonew`](https://pkg.go.dev/golang.org/x/tools/cmd/gonew)**
```bash
go install golang.org/x/tools/cmd/gonew@latest
gonew github.com/hylmithecoder/templatenextjs_golangserver github.com/you/new_project
```

`make rename` (`tools/rename`) updates `go.mod` and every quoted import path in the `.go` files, skipping `webui/`, `node_modules/`, `build/`, and `.git/`. It then runs `go build ./...` to check the result. You can run it again later; it always reads the current module path from `go.mod`.
