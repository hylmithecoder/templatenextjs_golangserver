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
```
