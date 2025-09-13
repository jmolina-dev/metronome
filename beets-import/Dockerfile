FROM docker.io/library/golang:1.25.1-alpine AS builder

WORKDIR /src
COPY go.mod main.go ./

RUN CGO_ENABLED=0 go build -o /beets-import .

FROM lscr.io/linuxserver/beets:latest

# Copy the compiled binary from the 'builder' stage into our final image.
COPY --from=builder /beets-import .

# The container will now run our custom script instead of its default command.
CMD ["beets-import"]
