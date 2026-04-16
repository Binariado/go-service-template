# Stage 1 — build
FROM golang:1.24 AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o service .

# Stage 2 — minimal runtime image
FROM gcr.io/distroless/static:nonroot

WORKDIR /app
COPY --from=builder /app/service .

EXPOSE $PORT

ENTRYPOINT ["/app/service"]
