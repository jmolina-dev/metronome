FROM golang:1.25.1-alpine AS builder

WORKDIR /src
COPY main.go .

RUN CGO_ENABLED=0 go build -o /beets-server .

FROM lscr.io/linuxserver/beets:latest

# Copy the compiled binary from the 'builder' stage into our final image.
COPY --from=builder /beets-server .

# The container will now run our custom script instead of its default command.
CMD ["beets-server"]
