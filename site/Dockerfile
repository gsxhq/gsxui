# site/Dockerfile — builds the gsxui showcase site for Cloud Run.
#
# Build context is the repo root (the Vite build reads ui/**/*.gsx and
# site/**/*.gsx for Tailwind's @source scan, and the Go module lives at the
# root) — from the repo root:
#
#   docker build -f site/Dockerfile -t gsxui-site .
#
# The .x.go files are committed (see Makefile's `generate`/`check`), so no
# `gsx generate` step runs in-container — the go stage does a plain
# `go build`. Verify this assumption holds before trusting a green build:
# `git status --porcelain -- '*.x.go'` must be clean on the commit being
# built (make check already enforces this pre-commit).

# ---- web stage: Vite production build (site/main.go embeds the result) ----
FROM node:22-bookworm-slim AS web
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npx vite build

# ---- go stage: compile the server against the committed .x.go + built dist ----
FROM golang:1.26.1-bookworm AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Overlay the freshly built assets — site/dist is .dockerignore'd from the
# context (it's a build artifact, never trust the host copy) so main.go's
# //go:embed all:dist has real content only after this COPY.
COPY --from=web /app/site/dist ./site/dist
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /gsxui-site ./site

# ---- run stage: distroless static, non-root ----
FROM gcr.io/distroless/static-debian12:nonroot AS run
COPY --from=build /gsxui-site /gsxui-site
# Cloud Run injects $PORT (defaults to 8080) and expects the container to
# bind it; site/main.go reads PORT with an 8080 fallback (GO_PORT still
# takes precedence, for parity with the local dev loop — see main.go).
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/gsxui-site"]
