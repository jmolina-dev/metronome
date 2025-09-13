FROM docker.io/library/golang:1.25.1-alpine AS builder

WORKDIR /src

COPY go.mod main.go ./

RUN CGO_ENABLED=0 go build -o /navidrome-scan .


# Stage 2: The "final" stage, where we create the minimal production image.
# We start from a fresh, tiny Alpine Linux image.
FROM gcr.io/distroless/static-debian12:nonroot

# We only copy the single, compiled binary from the 'builder' stage.
# None of the Go source code or build tools are included in the final image.
COPY --from=builder /navidrome-scan /navidrome-scan

# Set the entrypoint for the container. When the container starts, it will
# execute our compiled Go application.
ENTRYPOINT ["/navidrome-scan"]
