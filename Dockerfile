# Stage 1 — build
FROM golang:1.25 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make sqlc
RUN CGO_ENABLED=0 GOOS=linux go build -o service ./cmd/api

# Stage 2 — minimal runtime image (migrations are embedded in the binary)
FROM gcr.io/distroless/static:nonroot

WORKDIR /app
COPY --from=builder /app/service .

EXPOSE 3000

ENTRYPOINT ["/app/service"]
