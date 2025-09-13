# Metronome

A self-hosted Go webhook that orchestrates a music library workflow, running a Beets import and then automatically triggering a Navidrome library scan.

This service is designed to be the central automation point in a self-hosted music stack. It eliminates the need for manual beet import commands and polling-based cron jobs, creating a more efficient, event-driven pipeline.

### Features
- Webhook Trigger: Exposes a simple POST endpoint to start the entire task on demand.
- Automated Orchestration: After a successful Beets import, it automatically calls a separate Navidrome scan webhook to update your library immediately.
- Concurrency Safe: Uses a lock to ensure that only one import/scan task can run at a time, preventing data corruption.
- Built for Docker: Designed to run as a lightweight "sidecar" inside a custom Beets Docker image, communicating with other services over the container network.
- Configurable: All external endpoints and internal ports are configurable via environment variables.

### How It Works

This application is a small Go web server that runs alongside the standard Beets command-line tool in a single Docker container.
1. It listens for incoming POST requests on its /task/start endpoint.
2. Upon receiving a request, it runs the beet import command as a background process.
3. If the import is successful, it then makes a POST request to a configurable Navidrome scan webhook URL to complete the pipeline.
4. The service responds immediately with a 202 Accepted status, acknowledging that the task has been started.

### Usage
This service is intended to be built into a custom Docker image. See the accompanying Dockerfile and docker-compose.yml for a complete setup example.

#### Executing the task

Once the service is running, you can start the task by sending a POST request to its endpoint.

```curl -X POST http://localhost:9001/task/start```

This will kick off the entire process of importing new music with Beets and then updating the Navidrome library.
